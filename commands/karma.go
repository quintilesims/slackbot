package commands

import (
	"fmt"
	"io"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/urfave/cli"
)

func NewKarmaCommand(store db.Store, w io.Writer) cli.Command {
	return cli.Command{
		Name:      "!karma",
		Usage:     "display karma for the given key",
		ArgsUsage: "KEY",
		Action: func(c *cli.Context) error {
			key := c.Args().Get(0)
			if key == "" {
				return fmt.Errorf("Arg KEY is required")
			}

			karma := models.Karma{}
			if err := store.Read(models.StoreKeyKarma, &karma); err != nil {
				return err
			}

			text := fmt.Sprintf("karma for '%s': %d", key, karma[key])
			if _, err := w.Write([]byte(text)); err != nil {
				return err
			}

			return nil
		},
	}
}
