package bot

import (
	"fmt"
	"io"

	"github.com/quintilesims/slackbot/db"
	glob "github.com/ryanuber/go-glob"
	"github.com/urfave/cli"
)

// NewGlossaryCommand returns a cli.Command that manages !karma
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
		return fmt.Errorf("glossary define not implemented")
	}
}

func newGlossaryRemoveAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return fmt.Errorf("glossary remove not implemented")
	}
}

func newGlossarySearchAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		// example using glob matching
		glob.Glob("foo", "bar")

		return fmt.Errorf("glossary search not implemented")
	}
}
