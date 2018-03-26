package bot

import (
	"fmt"
	"io"
	"strings"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/urfave/cli"
)

// NewCandidateCommand returns a cli.Command that manages !candidate
func NewCandidateCommand(store db.Store, w io.Writer) cli.Command {
	return cli.Command{
		Name:  "!candidate",
		Usage: "manage candidates",
		Subcommands: []cli.Command{
			{
				Name:      "add",
				Usage:     "add a new candidate",
				ArgsUsage: "NAME",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "meta",
						Usage: "metadata about the candidate in key=val format",
					},
				},
				Action: newCandidateAddAction(store, w),
			},
			{
				Name:      "ls",
				Usage:     "list candidates",
				ArgsUsage: "[GLOB]",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "count",
						Value: 10,
						Usage: "The maximum number of candidates to display",
					},
					cli.BoolFlag{
						Name:  "descending",
						Usage: "Show results in descending alphabetical order",
					},
				},
				Action: newCandidateListAction(store, w),
			},
			{
				Name:      "rm",
				Usage:     "remove a candidate",
				ArgsUsage: "KEY",
				Action:    newCandidateRemoveAction(store, w),
			},
			{
				Name:      "update",
				Usage:     "update a candidate",
				ArgsUsage: "NAME",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "meta",
						Usage: "metadata about the candidate in key=val format",
					},
				},
				Action: newCandidateAddAction(store, w),
			},
		},
	}
}

func newCandidateAddAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		name := strings.Join(c.Args(), " ")
		if name == "" {
			return fmt.Errorf("NAME is required")
		}

		meta, err := parseMetaFlag(c.StringSlice("meta"))
		if err != nil {
			return err
		}

		candidates := models.Candidates{}
		if err := store.Read(db.CandidatesKey, &candidates); err != nil {
			return err
		}

		if _, ok := candidates[name]; ok {
			return fmt.Errorf("Candidate with name '%s' already exists", name)
		}

		candidates[name] = meta
		if err := store.Write(db.CandidatesKey, candidates); err != nil {
			return err
		}

		text := fmt.Sprintf("Ok, I've added a new candidate named *%s*", name)
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func newCandidateListAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return nil
	}
}

func newCandidateRemoveAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return nil
	}
}

func newCandidateUpdateAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return nil
	}
}

func parseMetaFlag(inputs []string) (map[string]string, error) {
	meta := map[string]string{}
	for _, input := range inputs {
		split := strings.Split(input, "=")
		if len(split) != 2 {
			return nil, fmt.Errorf("Input '%s' is not in key=val format", input)
		}

		meta[split[0]] = split[1]
	}

	return meta, nil
}
