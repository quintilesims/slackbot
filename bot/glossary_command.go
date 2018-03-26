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
				Name:      "add",
				Usage:     "add or set an entry in the glossary",
				ArgsUsage: "KEY DEFINITION",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "force",
						Usage: "overwrite existing definition",
					},
				},
				Action: newGlossaryAddAction(store, w),
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
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "count",
						Value: 10,
						Usage: "The maximum number of entries to display",
					},
				},
				Action: newGlossarySearchAction(store, w),
			},
		},
	}
}

func newGlossaryAddAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()

		key := args.First()
		if key == "" {
			return fmt.Errorf("Arg KEY is required")
		}

		definition := strings.Join(args.Tail(), " ")
		if definition == "" {
			return fmt.Errorf("Arg DEFINITION is required")
		}

		glossary := models.Glossary{}
		if err := store.Read(db.GlossaryKey, &glossary); err != nil {
			return err
		}

		if _, ok := glossary[key]; ok && !c.Bool("force") {
			return fmt.Errorf("An entry for *%s* already exists", key)
		}

		glossary[key] = definition
		if err := store.Write(db.GlossaryKey, glossary); err != nil {
			return err
		}

		text := fmt.Sprintf("OK, I've added *%s* as \"%s\"\n", key, definition)
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

		if _, ok := glossary[key]; !ok {
			return fmt.Errorf("There is no entry for *%s* in the glossary", key)
		}

		delete(glossary, key)
		if err := store.Write(db.GlossaryKey, glossary); err != nil {
			return err
		}

		text := fmt.Sprintf("Ok, I've deleted the entry for *%s*", key)
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

		definitions := models.Glossary{}
		for k, v := range glossary {
			if glob.Glob(g, k) {
				definitions[k] = v
			}
		}

		if len(definitions) == 0 {
			return fmt.Errorf("Could not find any glossary entries matching *%s*", g)
		}

		var text string
		keys := definitions.SortKeys(true)
		for i := 0; i < c.Int("count") && i < len(keys); i++ {
			key := keys[i]
			text += fmt.Sprintf("*%s*: %s\n", key, definitions[key])
		}

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}
