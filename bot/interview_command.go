package bot

import (
	"fmt"
	"io"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/urfave/cli"
)

// common time layouts
const (
	DateLayout = "01/02"
)

// NewInterviewCommand returns a cli.Command that manages !interview
func NewInterviewCommand(client SlackClient, store db.Store, w io.Writer) cli.Command {
	return cli.Command{
		Name:  "!interview",
		Usage: "manage interviews",
		Subcommands: []cli.Command{
			{
				Name:      "add",
				Usage:     "add a new interview",
				ArgsUsage: "NAME DATE (mm/dd)",
				Action:    newInterviewAddAction(client, store, w),
			},
		},
	}
}

func newInterviewAddAction(client SlackClient, store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		// todo: get user id
		name := c.Args().Get(0)
		if name == "" {
			return fmt.Errorf("Argument NAME is required")
		}

		date := c.Args().Get(1)
		if date == "" {
			return fmt.Errorf("Argument DATE is required")
		}

		t, err := time.Parse(DateLayout, date)
		if err != nil {
			return err
		}

		fmt.Printf("stuff: ", name, t)

		return fmt.Errorf("Not implemented")
	}
}

func newInterviewListAction(client SlackClient, store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return fmt.Errorf("Not implemented")
	}
}

func newInterviewRemoveAction(client SlackClient, store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return fmt.Errorf("Not implemented")
	}
}
