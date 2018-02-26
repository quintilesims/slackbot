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

				w := utils.WriterFunc(func(p []byte) (n int, err error) {
					text := string(p)
					log.Printf("[DEBUG] %s", text)
					msg := rtm.NewOutgoingMessage(text, d.Msg.Channel)
                                        rtm.SendMessage(msg)
                                        return len(p), nil
                                })

				handler.Writer = w
				handler.ErrWriter = w
				handler.Commands = []cli.Command{
					commands.NewEchoCommand(w).Command(),
				}

				if err := handler.Run(args); err != nil {
					text := fmt.Sprintf("Error: %s", err.Error())
					w.Write([]byte(text))
				}
			case *slack.RTMError:
				log.Printf("[ERROR] Unexected RTM error: %s", d.Msg)
			 case *slack.AckErrorEvent:
                                log.Printf("[ERROR] Unexpected Ack error: %s", d.Error())
			case *slack.InvalidAuthEvent:
				return fmt.Errorf("The bot's auth token is invalid")
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
