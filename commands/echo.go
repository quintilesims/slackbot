package commands

import (
	"fmt"
	"io"
	"strings"

	"github.com/urfave/cli"
)

type EchoCommand struct {
	w io.Writer
}

func NewEchoCommand(w io.Writer) *EchoCommand {
	return &EchoCommand{
		w: w,
	}
}

func (e EchoCommand) Command() cli.Command {
	return cli.Command{
		Name:      "echo",
		Usage:     "- todo -",
		ArgsUsage: "- todo - ",
		Action: func(c *cli.Context) error {
			return e.run(c.Args())
		},
	}
}

func (e *EchoCommand) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify at least one argument")
	}

	text := strings.Join(args, " ")
	if _, err := e.w.Write([]byte(text)); err != nil {
		return err
	}

	return nil
}
