package commands

import (
	"strings"

	"github.com/nlopes/slack"
	"github.com/urfave/cli"
)

type EchoCommand struct {
	rtm       *slack.RTM
	channelID string
}

func NewEchoCommand(rtm *slack.RTM, channelID string) *EchoCommand {
	return &EchoCommand{
		rtm:       rtm,
		channelID: channelID,
	}
}

func (e EchoCommand) Command() cli.Command {
	return cli.Command{
		Name:      "echo",
		Usage:     "- todo -",
		ArgsUsage: "- todo - ",
		Action:    e.echo,
	}
}

func (e *EchoCommand) echo(c *cli.Context) error {
	// cannot send empty messages in slack
	text := strings.Join(c.Args(), " ")
	if text == "" {
		text = "Hello, World!"
	}

	msg := e.rtm.NewOutgoingMessage(text, e.channelID)
	e.rtm.SendMessage(msg)
	return nil
}
