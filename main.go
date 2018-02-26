package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/commands"
	"github.com/quintilesims/slackbot/logging"
	"github.com/quintilesims/slackbot/utils"
	"github.com/urfave/cli"
)

var Version string

func main() {
	if Version == "" {
		Version = "unset/develop"
	}

	// todo: would be nice to wrap help text in markdown as code snippet
	// cli.AppHelpTemplate = fmt.Sprintf("```\n%s\n```", cli.AppHelpTemplate)

	app := cli.NewApp()
	app.Name = "slackbot"
	app.Version = Version
	app.Flags = []cli.Flag{
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

	app.Before = func(c *cli.Context) error {
		log.SetOutput(logging.NewLogWriter(c.Bool("debug")))
		logger := log.New(os.Stdout, "[DEBUG] ", log.Flags())
		slack.SetLogger(logger)

		return nil
	}

	app.Action = func(c *cli.Context) error {
		token := c.String("token")
		if token == "" {
			return fmt.Errorf("Token is not set!")
		}

		api := slack.New(token)
		rtm := api.NewRTM()
		defer rtm.Disconnect()

		handler := cli.NewApp()
		handler.Name = "slackbot"
		handler.Version = Version
		handler.ExitErrHandler = func(c *cli.Context, err error) {}

		// todo: how is this going to work in docker?
		go rtm.ManageConnection()
		for e := range rtm.IncomingEvents {
			switch d := e.Data.(type) {
			case *slack.ConnectedEvent:
				log.Printf("Slack connection successful!")
			case *slack.MessageEvent:
				args := strings.Split(d.Msg.Text, " ")
				if len(args) == 0 || args[0] != "slackbot" {
					continue
				}

				writeToSlack := func(text string) {
					msg := rtm.NewOutgoingMessage(text, d.Msg.Channel)
					rtm.SendMessage(msg)
				}

				handler.Writer = utils.WriterFunc(func(p []byte) (n int, err error) {
					writeToSlack(string(p))
					return len(p), nil
				})

				handler.ErrWriter = utils.WriterFunc(func(p []byte) (n int, err error) {
					log.Printf("[ERROR] %s", string(p))
					writeToSlack(string(p))
					return len(p), nil
				})

				handler.Commands = []cli.Command{
					commands.NewEchoCommand(rtm, d.Msg.Channel).Command(),
				}

				if err := handler.Run(args); err != nil {
					log.Printf("[ERROR] %s", err)
					writeToSlack(err.Error())
				}
			case *slack.RTMError:
				return fmt.Errorf("RTM Error: %s", d.Msg)
			case *slack.InvalidAuthEvent:
				return fmt.Errorf("The bot's auth token is invalid")
			case *slack.AckErrorEvent:
				return d
			default:
				log.Printf("[DEBUG] Unhandled event: %#v", e)
			}
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
