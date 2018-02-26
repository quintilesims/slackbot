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
		Action:    e.echo,
	}
}

func (e *EchoCommand) echo(c *cli.Context) error {
	text := strings.Join(c.Args(), " ")
	if text == "" {
		return fmt.Errorf("please specify at least one argument")
	}

	if _, err := e.w.Write([]byte(text)); err != nil {
		return err
	}

	return nil
}
