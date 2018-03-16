package bot

import (
	"fmt"
	"io"
	"strings"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	glob "github.com/ryanuber/go-glob"
	"github.com/urfave/cli"
)

// NewGlossaryCommand returns a cli.Command that manages !glossary
func NewGlossaryCommand(store db.Store, w io.Writer) cli.Command {
	return cli.Command{
		Name:  "!glossary",
		Usage: "manage the glossary",
		Subcommands: []cli.Command{
			{
				Name:      "define",
				Usage:     "add or set an entry in the glossary",
				ArgsUsage: "KEY DEFINITION",
				Action:    newGlossaryDefineAction(store, w),
			},
			{
				Name:      "rm",
				Usage:     "remove an entry from the glossary",
				ArgsUsage: "KEY",
				Action:    newGlossaryRemoveAction(store, w),
			},
			{
				Name:      "search",
				Usage:     "search for entries in the glossary",
				ArgsUsage: "GLOB",
				Action:    newGlossarySearchAction(store, w),
			},
		},
	}
}

func newGlossaryDefineAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()
		key := args.First()
		definition := strings.Join(args.Tail(), " ")

		if key == "" {
			return fmt.Errorf("Arg KEY is required")
		}

		if definition == "" {
			return fmt.Errorf("Arg DEFINITION is required")
		}

		glossary := models.Glossary{}
		if err := store.Read(db.GlossaryKey, &glossary); err != nil {
			return err
		}

		var text string
		if previousDefinition, ok := glossary[key]; previousDefinition != "" && ok {
			text = fmt.Sprintf("Replacing previous definition of *%s*: *%s* with *%s*", key, previousDefinition, definition)
		} else {
			text = fmt.Sprintf("Ok, *%s* %s", key, definition)
		}

		glossary[key] = definition

		if err := store.Write(db.GlossaryKey, glossary); err != nil {
			return err
		}

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func newGlossaryRemoveAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		key := c.Args().Get(0)

		if key == "" {
			return fmt.Errorf("Arg KEY is required")
		}

		glossary := models.Glossary{}
		if err := store.Read(db.GlossaryKey, &glossary); err != nil {
			return err
		}

		if _, ok := glossary[key]; ok {
			delete(glossary, key)
		} else {
			return fmt.Errorf("Key: *%s* not in glossary", key)
		}

		if err := store.Write(db.GlossaryKey, glossary); err != nil {
			return err
		}

		text := fmt.Sprintf("Ok, deleted *%s*", key)
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func newGlossarySearchAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		g := c.Args().Get(0)

		if g == "" {
			return fmt.Errorf("Arg GLOB is required")
		}

		glossary := models.Glossary{}
		if err := store.Read(db.GlossaryKey, &glossary); err != nil {
			return err
		}

		results := models.Glossary{}
		for k, v := range glossary {
			if glob.Glob(g, k) {
				results[k] = v
			}
		}

		if len(results) == 0 {
			return fmt.Errorf("Could not find any glossary entries matching the specified pattern *%s*", g)
		}

		var text string
		for key, definition := range results {
			text += fmt.Sprintf("*%s* %s\n", key, definition)
		}

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}
