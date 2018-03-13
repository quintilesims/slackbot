package bot

import (
	"io"

	"github.com/urfave/cli"
)

// NewPingCommand returns a cli.Command that manages !ping
func NewPingCommand(w io.Writer) cli.Command {
	return cli.Command{
		Name:      "!ping",
		Usage:     "ping the slackbot",
		ArgsUsage: " ",
		Action: func(c *cli.Context) error {
			if _, err := w.Write([]byte("pong")); err != nil {
				return err
			}

			return nil
		},
	}
}
