package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/behaviors"
	"github.com/quintilesims/slackbot/commands"
	"github.com/quintilesims/slackbot/common"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/utils"
	"github.com/urfave/cli"
)

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
			Name:   "token",
			Usage:  "authentication token for the slack bot",
			EnvVar: "SB_TOKEN",
		},
	}

	var api *slack.Client
	var store db.Store
	var behavs []behaviors.Behavior

	slackbot.Before = func(c *cli.Context) error {
		debug := c.Bool("debug")
		log.SetOutput(utils.NewLogWriter(debug))

		token := c.String("token")
		if token == "" {
			return fmt.Errorf("Token is not set!")
		}

		api = slack.New(token)
		api.SetDebug(debug)

		store = db.NewMemoryStore()
		if err := common.Init(store); err != nil {
			return err
		}

		behavs = []behaviors.Behavior{
			behaviors.NewKarmaTrackingBehavior(store),
		}

		return nil
	}

	slackbot.Action = func(c *cli.Context) error {
		rtm := api.NewRTM()
		defer rtm.Disconnect()

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

				eventApp.Commands = []cli.Command{
					commands.NewEchoCommand(w),
					commands.NewKarmaCommand(store, w),
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
