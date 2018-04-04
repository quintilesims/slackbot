package bot

import (
	"fmt"
	"io"
	"strings"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/urfave/cli"
)

// NewAliasCommand returns a cli.Command that manages !alias
func NewAliasCommand(store db.Store, w io.Writer, invalidate func()) cli.Command {
	return cli.Command{
		Name:  "!alias",
		Usage: "manage aliases",
		Subcommands: []cli.Command{
			{
				Name:  "add",
				Usage: "add a new alias",
				// todo: give better help text
				ArgsUsage: "NAME VALUE",
				Action:    newAliasAddAction(store, w, invalidate),
			},
			{
				Name:      "ls",
				Usage:     "list all aliases",
				ArgsUsage: " ",
				Action:    newAliasListAction(store, w),
			},
			{
				Name:      "rm",
				Usage:     "remove an alias",
				ArgsUsage: "NAME",
				Action:    newAliasRemoveAction(store, w, invalidate),
			},
			{
				Name:      "test",
				Usage:     "test an alias",
				ArgsUsage: "TEXT",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "channel",
						Value: "channel_id",
						Usage: "the channel id to use for the test message",
					},
					cli.StringFlag{
						Name:  "user",
						Value: "user_id",
						Usage: "the user id to use for the test message",
					},
				},
				Action: newAliasTestAction(store, w),
			},
		},
	}
}

func newAliasAddAction(store db.Store, w io.Writer, invalidate func()) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()
		name := args.Get(0)
		if name == "" {
			return fmt.Errorf("Argument NAME is required")
		}

		if strings.Contains(name, " ") {
			return fmt.Errorf("Alias names may not contain whitespace")
		}

		value := strings.Join(args[1:], " ")
		if value == "" {
			return fmt.Errorf("Argument VALUE is required")
		}

		aliases := models.Aliases{}
		if err := store.Read(db.AliasesKey, &aliases); err != nil {
			return err
		}

		if _, ok := aliases[name]; ok {
			return fmt.Errorf("An alias for *%s* already exists", name)
		}

		aliases[name] = value
		if err := store.Write(db.AliasesKey, aliases); err != nil {
			return err
		}

		// make sure we tell the alias behavior cache to invalidate
		invalidate()

		return writef(w, "Ok, I've added a new alias for *%s*", name)
	}
}

func newAliasListAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		aliases := models.Aliases{}
		if err := store.Read(db.AliasesKey, &aliases); err != nil {
			return err
		}

		if len(aliases) == 0 {
			return fmt.Errorf("I don't have any aliases at the moment")
		}

		text := "Here are the aliases I have: \n"
		for name, value := range aliases {
			text += fmt.Sprintf("*%s*: `%s`\n", name, value)
		}

		return write(w, text)
	}
}

func newAliasRemoveAction(store db.Store, w io.Writer, invalidate func()) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		name := c.Args().Get(0)
		if name == "" {
			return fmt.Errorf("Argument NAME is required")
		}

		aliases := models.Aliases{}
		if err := store.Read(db.AliasesKey, &aliases); err != nil {
			return err
		}

		if _, ok := aliases[name]; !ok {
			return fmt.Errorf("No aliases for *%s* exist", name)
		}

		delete(aliases, name)
		if err := store.Write(db.AliasesKey, aliases); err != nil {
			return err
		}

		// make sure we tell the alias behavior cache to invalidate
		invalidate()

		return writef(w, "Ok, I've removed alias *%s*", name)
	}
}

func newAliasTestAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		text := strings.Join(c.Args(), " ")
		if text == "" {
			return fmt.Errorf("Argument TEXT is required")
		}

		aliases := models.Aliases{}
		if err := store.Read(db.AliasesKey, &aliases); err != nil {
			return err
		}

		m := &slack.MessageEvent{
			Msg: slack.Msg{
				Channel: c.String("channel"),
				User:    c.String("user"),
				Text:    text,
			},
		}

		if err := aliases.Apply(m); err != nil {
			return err
		}

		return write(w, m.Text)
	}
}
