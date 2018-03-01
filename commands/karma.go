package commands

import (
	"fmt"
	"io"

	"github.com/quintilesims/slackbot/common"
	"github.com/quintilesims/slackbot/db"
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

			karma, err := readKarma(key, store)
			if err != nil {
				return err
			}

			text := fmt.Sprintf("karma for '%s': %d", key, karma)
			if _, err := w.Write([]byte(text)); err != nil {
				return err
			}

			return nil
		},
	}
}

func readKarma(key string, store db.Store) (int, error) {
	karma := map[string]int{}
	if err := store.Read(common.StoreKeyKarma, &karma); err != nil {
		return 0, err
	}

	return karma[key], nil
}
