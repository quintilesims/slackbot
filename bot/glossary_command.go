package bot

import (
	"fmt"
	"io"
	"strings"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
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
				Usage:     "add an entry in the glossary",
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
				Name:  "ls",
				Usage: "list entries in the glossary",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "count",
						Value: 10,
						Usage: "The maximum number of entries to display",
					},
				},
				Action: newGlossaryListAction(store, w),
			},
			{
				Name:      "rm",
				Usage:     "remove an entry from the glossary",
				ArgsUsage: "KEY",
				Action:    newGlossaryRemoveAction(store, w),
			},
			{
				Name:      "show",
				Usage:     "show an entry in the glossary",
				ArgsUsage: "ENTRY",
				Action:    newGlossaryShowAction(store, w),
			},
		},
	}
}

func newGlossaryAddAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()

		key := args.First()
		if key == "" {
			return fmt.Errorf("Argument KEY is required")
		}

		definition := strings.Join(args.Tail(), " ")
		if definition == "" {
			return fmt.Errorf("Argument DEFINITION is required")
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

		return writef(w, "OK, I've added *%s* as \"%s\"\n", key, definition)
	}
}

func newGlossaryListAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		glossary := models.Glossary{}
		if err := store.Read(db.GlossaryKey, &glossary); err != nil {
			return err
		}

		if len(glossary) == 0 {
			return fmt.Errorf("There are currently no entries in the glossary")
		}

		var text string
		keys := glossary.SortKeys(true)
		for i := 0; i < len(keys) && i < c.Int("count"); i++ {
			entry := keys[i]
			text += fmt.Sprintf("*%s*: %s\n", entry, glossary[entry])
		}

		return write(w, text)
	}
}

func newGlossaryRemoveAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		key := strings.Join(c.Args(), " ")
		if key == "" {
			return fmt.Errorf("Argument KEY is required")
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

		return writef(w, "Ok, I've deleted the entry for *%s*", key)
	}
}

func newGlossaryShowAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		entry := strings.ToLower(strings.Join(c.Args(), " "))
		if entry == "" {
			return fmt.Errorf("ENTRY is required")
		}

		glossary := models.Glossary{}
		if err := store.Read(db.GlossaryKey, &glossary); err != nil {
			return err
		}

		var definition string
		for k, v := range glossary {
			if strings.ToLower(k) == entry {
				definition = v
				break
			}
		}

		if definition == "" {
			return fmt.Errorf("There is no entry for *%s*", entry)
		}

		return writef(w, "*%s*: %s\n", entry, definition)
	}
}
