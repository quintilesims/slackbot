package bot

import (
	"io"
	"strings"

	"github.com/urfave/cli"
)

// NewEchoCommand returns a cli.Command that manages !echo
func NewEchoCommand(w io.Writer) cli.Command {
	return cli.Command{
		Name:            "!echo",
		Usage:           "display the given message",
		ArgsUsage:       "[args...]",
		SkipFlagParsing: true,
		Action: func(c *cli.Context) error {
			text := strings.Join(c.Args(), " ")
			return write(w, text)
		},
	}
}
