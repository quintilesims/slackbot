package bot

import (
	"fmt"
	"io"

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
				ArgsUsage: "NAME TEXT",
				Action:    newAliasTestAction(store, w),
			},
		},
	}
}

func newAliasAddAction(store db.Store, w io.Writer, invalidate func()) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()
		name := args.Get(0)
		if name == "" {
			return fmt.Errorf("NAME is required")
		}

		value := args.Get(1)
		if value == "" {
			return fmt.Errorf("VALUE is required")
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

		text := fmt.Sprintf("Ok, I've added a new alias for *%s*", name)
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
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

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func newAliasRemoveAction(store db.Store, w io.Writer, invalidate func()) func(c *cli.Context) error {
	return func(c *cli.Context) error {

		// make sure we tell the alias behavior cache to invalidate
		invalidate()

		return nil
	}
}

func newAliasTestAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return nil
	}
}
