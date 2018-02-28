package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/logging"
	"github.com/quintilesims/slackbot/rtm"
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
			Name:  "d, debug",
			Usage: "- todo -",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "- todo -",
			EnvVar: "SB_TOKEN",
		},
	}

	slackbot.Before = func(c *cli.Context) error {
		debug := c.Bool("debug")
		log.SetOutput(logging.NewLogWriter(debug))

		return nil
	}

	slackbot.Action = func(c *cli.Context) error {
		token := c.String("token")
		if token == "" {
			return fmt.Errorf("Token is not set!")
		}

		store := db.NewMemoryStore()
		api := slack.New(token)
		api.SetDebug(c.Bool("debug"))

		behaviors := rtm.Behaviors{
			rtm.NewEchoBehavior(),
			rtm.NewKarmaBehavior(store),
		}

		behaviors = append(behaviors, rtm.NewHelpBehavior(behaviors...))
		if err := behaviors.Init(); err != nil {
			return err
		}

		r := api.NewRTM()
		defer r.Disconnect()

		newChannelWriter := func(channelID string) io.Writer {
			return utils.WriterFunc(func(b []byte) (n int, err error) {
				msg := r.NewOutgoingMessage(string(b), channelID)
				r.SendMessage(msg)
				return len(b), nil
			})
		}

		go r.ManageConnection()
		for e := range r.IncomingEvents {
			switch event := e.Data.(type) {
			case *slack.ConnectedEvent:
				log.Printf("[INFO] Slack connection successful!")
			case *slack.MessageEvent:
				w := newChannelWriter(event.Msg.Channel)
				if err := behaviors.OnMessageEvent(event, w); err != nil {
					w.Write([]byte(err.Error()))
				}
			case *slack.RTMError:
				log.Printf("[ERROR] Unexected RTM error: %s", event.Msg)
			case *slack.AckErrorEvent:
				log.Printf("[ERROR] Unexpected Ack error: %s", event.Error())
			case *slack.InvalidAuthEvent:
				return fmt.Errorf("The bot's auth token is invalid")
			default:
				log.Printf("[DEBUG] Unhandled event: %#v", e)
			}
		}

		return nil
	}

	if err := slackbot.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
