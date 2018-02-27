package rtm

import (
	"fmt"
	"io"
	"strings"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
)

const StoreKeyKarma = "karma"

func NewKarmaAction(store db.Store) *ActionSchema {
	return &ActionSchema{
		Name:  "karma",
		Usage: "!karma `id`",
		Help:  "Show karma for the specified ID.\nAdd or subtract karma by typing `++` or `--` after the ID (e.g. `dogs++`).",
		Init: func() error {
			var karma map[string]int
			if err := store.Read(StoreKeyKarma, &karma); err != nil {
				if _, ok := err.(db.MissingEntryError); ok {
					return store.Write(StoreKeyKarma, karma)
				}

				return err
			}

			return nil
		},
		OnMessageEvent: func(e *slack.MessageEvent, w io.Writer) error {
			args := strings.Split(e.Msg.Text, " ")
			if len(args) == 0 {
				return nil
			}

			if len(args) == 1 {
				var update func(i int) int

				switch {
				case args[0] == "!karma":
					return fmt.Errorf("!karma requires an identifier")
				case strings.HasSuffix(args[0], "++"):
					update = func(i int) int { return i + 1 }
				case strings.HasSuffix(args[0], "--"):
					update = func(i int) int { return i - 1 }
				default:
					return nil
				}

				karma := map[string]int{}
				if err := store.Read(StoreKeyKarma, &karma); err != nil {
					return err
				}

				// strip the '++' or '--' from the id
				id := args[0][:len(args[0])-2]
				karma[id] = update(karma[id])
				return store.Write(StoreKeyKarma, karma)
			}

			if args[0] == "!karma" {
				if len(args) > 2 {
					return fmt.Errorf("!karma can only accept a single identifier")
				}

				karma := map[string]int{}
				if err := store.Read(StoreKeyKarma, &karma); err != nil {
					return err
				}

				id := args[1]
				text := fmt.Sprintf("karma for '%s': %d", id, karma[id])
				if _, err := w.Write([]byte(text)); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
