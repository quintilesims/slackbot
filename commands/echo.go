package commands

import (
	"io"
	"strings"

	"github.com/urfave/cli"
)

// NewEchoCommand returns a cli.Command that manages !echo 
func NewEchoCommand(w io.Writer) cli.Command {
	return cli.Command{
		Name:      "!echo",
		Usage:     "display the given message",
		ArgsUsage: "[args...]",
		Action: func(c *cli.Context) error {
			text := strings.Join(c.Args(), " ")
			if _, err := w.Write([]byte(text)); err != nil {
				return err
			}

			return nil
		},
	}
}
