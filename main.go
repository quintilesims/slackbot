package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/bot"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/runner"
	"github.com/quintilesims/slackbot/utils"
	"github.com/urfave/cli"
)

// Version of the application
var Version string

func main() {
	if Version == "" {
		Version = "unset/develop"
	}

	slackbot := cli.NewApp()
	slackbot.Name = "slackbot"
	slackbot.Version = Version
	slackbot.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "d, debug",
			Usage:  "enable debug logging",
			EnvVar: "SB_DEBUG",
		},
		cli.StringFlag{
			Name:   "slack-bot-token",
			Usage:  "authentication token for the slack bot",
			EnvVar: "SB_SLACK_BOT_TOKEN",
		},
		cli.StringFlag{
			Name:   "slack-app-token",
			Usage:  "authentication token for the slack application",
			EnvVar: "SB_SLACK_APP_TOKEN",
		},
		cli.StringFlag{
			Name:   "tenor-key",
			Usage:  "authentication key for the Tenor API",
			EnvVar: "SB_TENOR_KEY",
		},
		cli.StringFlag{
			Name:   "aws-region",
			Usage:  "region for aws api",
			Value:  "us-west-2",
			EnvVar: "SB_AWS_REGION",
		},
		cli.StringFlag{
			Name:   "aws-access-key",
			Usage:  "access key for aws api",
			EnvVar: "SB_AWS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "aws-secret-key",
			Usage:  "secret key for aws api",
			EnvVar: "SB_AWS_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "dynamodb-table",
			Usage:  "name of the dynamodb table",
			EnvVar: "SB_DYNAMODB_TABLE",
		},
	}

	var appClient *slack.Client
	var botClient *slack.Client
	var store db.Store

	slackbot.Before = func(c *cli.Context) error {
		rand.Seed(time.Now().UnixNano())

		debug := c.Bool("debug")
		log.SetOutput(utils.NewLogWriter(debug))

		botToken := c.String("slack-bot-token")
		if botToken == "" {
			return fmt.Errorf("Bot Token is not set!")
		}

		botClient = slack.New(botToken)
		botClient.SetDebug(debug)

		appToken := c.String("slack-app-token")
		if appToken == "" {
			return fmt.Errorf("App Token is not set!")
		}

		appClient = slack.New(appToken)
		appClient.SetDebug(debug)

		accessKey := c.String("aws-access-key")
		if accessKey == "" {
			return fmt.Errorf("AWS Access Key is not set!")
		}

		secretKey := c.String("aws-secret-key")
		if secretKey == "" {
			return fmt.Errorf("AWS Secret Key is not set!")
		}

		region := c.String("aws-region")
		if region == "" {
			return fmt.Errorf("AWS Region is not set!")
		}

		table := c.String("dynamodb-table")
		if table == "" {
			return fmt.Errorf("DynamoDB Table is not set!")
		}

		config := &aws.Config{
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
			Region:      aws.String(region),
		}

		store = db.NewDynamoDBStore(session.New(config), table)
		if err := db.Init(store); err != nil {
			return err
		}

		return nil
	}

	slackbot.Action = func(c *cli.Context) error {
		aliasBehavior := bot.NewAliasBehavior(store)
		behaviors := bot.Behaviors{
			bot.NewNormalizeTextBehavior(),
			aliasBehavior.Behavior(),
			bot.NewKarmaTrackingBehavior(store),
		}

		defer runner.NewCleanupRunner(store).RunEvery(time.Hour).Stop()
		defer runner.NewReminderRunner(store, botClient).RunEvery(time.Minute * 5).Stop()

		// initiate the RTM websocket connection
		rtm := botClient.NewRTM()
		defer rtm.Disconnect()
		go rtm.ManageConnection()

		// allow the RedoBehavior to send events to the rtm.IncomingEvents channel
		redoBehavior := bot.NewRedoBehavior(rtm.IncomingEvents)

		for event := range rtm.IncomingEvents {
			if err := behaviors.Run(event); err != nil {
				log.Printf("[ERROR] %v", err)
			}

			switch e := event.Data.(type) {
			case *slack.ConnectedEvent:
				log.Printf("[INFO] Slack connection successful!")
			case *slack.MessageEvent:
				if !strings.HasPrefix(e.Msg.Text, "!") {
					continue
				}

				args, err := utils.ParseShell("slackbot " + e.Msg.Text)
				if err != nil {
					msg := rtm.NewOutgoingMessage(err.Error(), e.Channel)
					rtm.SendMessage(msg)
					continue
				}

				var isDisplayingHelp bool
				w := bytes.NewBuffer(nil)

				app := cli.NewApp()
				app.Name = "slackbot"
				app.Usage = "making email obsolete one step at a time"
				app.UsageText = "command [flags...] arguments..."
				app.Version = Version
				app.Writer = utils.WriterFunc(func(b []byte) (n int, err error) {
					isDisplayingHelp = true
					return w.Write(b)
				})

				app.CommandNotFound = func(c *cli.Context, command string) {
					text := fmt.Sprintf("Command '%s' does not exist", command)
					w.WriteString(text)
				}

				app.Commands = []cli.Command{
					bot.NewAliasCommand(store, w, aliasBehavior.Invalidate),
					bot.NewCandidateCommand(store, w),
					bot.NewEchoCommand(w),
					bot.NewGIFCommand(bot.TenorAPIEndpoint, c.String("tenor-key"), w),
					bot.NewGlossaryCommand(store, w),
					bot.NewHelpCommand(w),
					bot.NewInterviewCommand(store, w),
					bot.NewKarmaCommand(store, w),
					bot.NewRedoCommand(func() error { return redoBehavior.Trigger(e.Msg.Channel) }),
					bot.NewTriviaCommand(store, bot.TriviaAPIEndpoint, w),
					bot.NewUndoCommand(appClient, botClient, e.Channel, rtm.GetInfo().User.ID),
				}

				// record the 'last event' for each channel
				redoBehavior.Record(e.Msg.Channel, event)

				if err := app.Run(args); err != nil {
					log.Printf("[ERROR] %v", err)
					w.WriteString(err.Error())
				}

				text := w.String()
				if isDisplayingHelp {
					text = fmt.Sprintf("```%s```", text)
				}

				msg := rtm.NewOutgoingMessage(text, e.Channel)
				rtm.SendMessage(msg)
			case *slack.InvalidAuthEvent:
				return fmt.Errorf("The bot's auth token is invalid")
			default:
				log.Printf("[DEBUG] Unhandled event: %#v", event)
			}
		}

		return nil
	}

	if err := slackbot.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
