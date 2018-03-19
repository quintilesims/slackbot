package bot

import (
	"fmt"
	"io"
	"sort"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/urfave/cli"
)

// NewAliasCommand returns a cli.Command that manages !alias
func NewAliasCommand(store db.Store, w io.Writer) cli.Command {
	return cli.Command{
		Name:  "!alias",
		Usage: "manage aliases",
		Subcommands: []cli.Command{
			{
				Name:  "add",
				Usage: "add a new alias",
				// todo: give better help text
				ArgsUsage: "NAME `PATTERN` `TEMPLATE`",
				Action:    newAliasAddAction(store, w),
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
				Action:    newAliasRemoveAction(store, w),
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

func newAliasAddAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()
		name := args.Get(0)
		if name == "" {
			return fmt.Errorf("NAME is required")
		}

		pattern := args.Get(1)
		if pattern == "" {
			return fmt.Errorf("PATTERN is required")
		}

		template := args.Get(2)
		if template == "" {
			return fmt.Errorf("TEMPLATE is required")
		}

		aliases := models.Aliases{}
		if err := store.Read(db.AliasesKey, &aliases); err != nil {
			return err
		}

		aliases[name] = models.Alias{
			Pattern:  pattern,
			Template: template,
		}

		if err := store.Write(db.AliasesKey, aliases); err != nil {
			return err
		}

		text := fmt.Sprintf("Ok, I've added a new alias named *%s*", name)
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

// todo: !alias show NAME that shows the pattern and template
func newAliasListAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		aliases := models.Aliases{}
		if err := store.Read(db.AliasesKey, &aliases); err != nil {
			return err
		}

		if len(aliases) == 0 {
			return fmt.Errorf("I don't have any aliases at the moment")
		}

		names := make([]string, 0, len(aliases))
		for name := range aliases {
			names = append(names, name)
		}

		sort.Sort(sort.StringSlice(names))

		text := "Here are the aliases I have: "
		for i, name := range names {
			if i == len(names)-1 {
				text += fmt.Sprintf("and *%s*", name)
				break
			}

			text += fmt.Sprintf("*%s*, ", name)
		}

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func newAliasRemoveAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return nil
	}
}

func newAliasTestAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return nil
	}
}
