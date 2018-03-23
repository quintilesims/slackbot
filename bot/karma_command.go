package bot

import (
	"fmt"
	"io"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/ryanuber/go-glob"
	"github.com/urfave/cli"
)

// NewKarmaCommand returns a cli.Command that manages !karma
func NewKarmaCommand(store db.Store, w io.Writer) cli.Command {
	return cli.Command{
		Name:      "!karma",
		Usage:     "display karma for entries that match the given GLOB",
		ArgsUsage: "GLOB",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "count",
				Value: 10,
				Usage: "The maximum number of entries to display",
			},
			cli.BoolFlag{
				Name:  "ascending",
				Usage: "Show results in ascending order",
			},
		},
		Action: func(c *cli.Context) error {
			g := c.Args().Get(0)
			if g == "" {
				return fmt.Errorf("Arg GLOB is required")
			}

			karmas := models.Karmas{}
			if err := store.Read(db.KarmasKey, &karmas); err != nil {
				return err
			}

			results := models.Karmas{}
			for k, v := range karmas {
				if glob.Glob(g, k) {
					results[k] = v
				}
			}

			if len(results) == 0 {
				return fmt.Errorf("Could not find any karma entries matching the specified pattern")
			}

			keys := results.SortKeys(c.Bool("ascending"))
			for i := 0; i < c.Int("count") && i < len(keys); i++ {
				karma := results[keys[i]]
				text := fmt.Sprintf("karma for *%s*: %s\n", keys[i], karma)
				if _, err := w.Write([]byte(text)); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
