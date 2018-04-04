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
				ArgsUsage: " ",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "count",
						Value: 10,
						Usage: "The maximum number of candidates to display",
					},
				},
				Action: newCandidateListAction(store, w),
			},
			{
				Name:      "rm",
				Usage:     "remove a candidate",
				ArgsUsage: "NAME",
				Action:    newCandidateRemoveAction(store, w),
			},
			{
				Name:      "info",
				Usage:     "show information about a candidate",
				ArgsUsage: "NAME",
				Action:    newCandidateInfoAction(store, w),
			},
			{
				Name:      "update",
				Usage:     "upsert a candidate's information",
				ArgsUsage: "NAME KEY VAL",
				Action:    newCandidateUpdateAction(store, w),
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

		return writef(w, "Ok, I've added a new candidate named *%s*", name)
	}
}

func newCandidateListAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		candidates := models.Candidates{}
		if err := store.Read(db.CandidatesKey, &candidates); err != nil {
			return err
		}

		if len(candidates) == 0 {
			return fmt.Errorf("I don't have any candidates at the moment")
		}

		text := "Here are the candidates I have: \n"
		keys := candidates.SortKeys(true)
		for i := 0; i < c.Int("count") && i < len(keys); i++ {
			text += fmt.Sprintf("*%s* \n", keys[i])
		}

		return write(w, text)
	}
}

func newCandidateRemoveAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		name := strings.Join(c.Args(), " ")
		if name == "" {
			return fmt.Errorf("NAME is required")
		}

		candidates := models.Candidates{}
		if err := store.Read(db.CandidatesKey, &candidates); err != nil {
			return err
		}

		if _, ok := candidates[name]; !ok {
			return fmt.Errorf("I don't have any candidates by the name *%s*", name)
		}

		delete(candidates, name)
		if err := store.Write(db.CandidatesKey, candidates); err != nil {
			return err
		}

		return writef(w, "Ok, I've deleted candidate *%s*", name)
	}
}

func newCandidateInfoAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		name := strings.Join(c.Args(), " ")
		if name == "" {
			return fmt.Errorf("NAME is required")
		}

		candidates := models.Candidates{}
		if err := store.Read(db.CandidatesKey, &candidates); err != nil {
			return err
		}

		metadata, ok := candidates[name]
		if !ok {
			return fmt.Errorf("I don't have any candidates by the name *%s*", name)
		}

		text := fmt.Sprintf("*%s*: \n", name)
		if len(metadata) == 0 {
			text += "This candidate currently doesn't have any information associated with them"
		}

		for key, val := range metadata {
			text += fmt.Sprintf("%s: %s\n", key, val)
		}

		return write(w, text)
	}
}

func newCandidateUpdateAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()
		name := args.Get(0)
		if name == "" {
			return fmt.Errorf("NAME is required")
		}

		key := args.Get(1)
		if key == "" {
			return fmt.Errorf("KEY is required")
		}

		val := args.Get(2)
		if val == "" {
			return fmt.Errorf("VAL is required")
		}

		candidates := models.Candidates{}
		if err := store.Read(db.CandidatesKey, &candidates); err != nil {
			return err
		}

		if _, ok := candidates[name]; !ok {
			return fmt.Errorf("I don't have any candidates by the name *%s*", name)
		}

		if candidates[name] == nil {
			candidates[name] = map[string]string{}
		}

		candidates[name][key] = val
		if err := store.Write(db.CandidatesKey, candidates); err != nil {
			return err
		}

		return writef(w, "Ok, I've updated information for *%s*", name)
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
