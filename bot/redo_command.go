package bot

import (
	"github.com/urfave/cli"
)

// NewRedoCommand returns a cli.Command that manages !redo.
func NewRedoCommand(trigger func() error) cli.Command {
	return cli.Command{
		Name:  "!redo",
		Usage: "redo the last command executed by the slackbot",
		Action: func(c *cli.Context) error {
			return trigger()
		},
	}
}
