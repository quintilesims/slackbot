package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/behaviors"
	"github.com/quintilesims/slackbot/commands"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/lock"
	"github.com/quintilesims/slackbot/runners"
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
			Name:   "slack-token",
			Usage:  "authentication token for the slack bot",
			EnvVar: "SB_SLACK_TOKEN",
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

	var client *slack.Client
	var store db.Store
	var behavs []behaviors.Behavior

	slackbot.Before = func(c *cli.Context) error {
		debug := c.Bool("debug")
		log.SetOutput(utils.NewLogWriter(debug))

		token := c.String("slack-token")
		if token == "" {
			return fmt.Errorf("Token is not set!")
		}

		client = slack.New(token)
		client.SetDebug(debug)

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

		behavs = []behaviors.Behavior{
			behaviors.NewKarmaTrackingBehavior(store),
		}

		return nil
	}

	slackbot.Action = func(c *cli.Context) error {
		rtm := client.NewRTM()
		defer rtm.Disconnect()

		remindersRunner := runners.NewRemindersRunner(lock.NewStoreLock("reminders", store), store, &rtm.Client)
		ticker := remindersRunner.RunEvery(time.Minute)
		defer ticker.Stop()

		newChannelWriter := func(channelID string) io.Writer {
			return utils.WriterFunc(func(b []byte) (n int, err error) {
				msg := rtm.NewOutgoingMessage(string(b), channelID)
				rtm.SendMessage(msg)
				return len(b), nil
			})
		}

		go rtm.ManageConnection()
		for event := range rtm.IncomingEvents {
			for _, b := range behavs {
				if err := b(event); err != nil {
					log.Printf("[ERROR] %v", err)
				}
			}

			switch e := event.Data.(type) {
			case *slack.ConnectedEvent:
				log.Printf("[INFO] Slack connection successful!")
			case *slack.MessageEvent:
				if !strings.HasPrefix(e.Msg.Text, "!") {
					continue
				}

				w := newChannelWriter(e.Msg.Channel)
				eventApp := cli.NewApp()
				eventApp.Writer = w
				eventApp.CommandNotFound = func(c *cli.Context, command string) {
					text := fmt.Sprintf("Command '%s' does not exist", command)
					w.Write([]byte(text))
				}

				generateID := utils.NewGUIDGenerator()
				userParser := utils.NewSlackUserParser(&rtm.Client)
				eventApp.Commands = []cli.Command{
					commands.NewEchoCommand(w),
					commands.NewKarmaCommand(store, w),
					commands.NewRemindersCommand(store, w, generateID, userParser),
				}

				args := append([]string{""}, utils.ParseShell(e.Msg.Text)...)
				if err := eventApp.Run(args); err != nil {
					w.Write([]byte(err.Error()))
				}
			case *slack.RTMError:
				log.Printf("[ERROR] Unexected RTM error: %s", e.Msg)
			case *slack.AckErrorEvent:
				if e.Error() == "Code 2 - message text is missing" {
					continue
				}

				log.Printf("[ERROR] Unexpected Ack error: %s", e.Error())
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
