package rtm

import (
	"fmt"
	"io"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
)

const StoreKeyKarma = "karma"

func NewKarmaAction(store db.Store) *ActionSchema {
	return &ActionSchema{
		Name:  "karma",
		Usage: "`@username++` or `@username--`",
		Help:  "add or subtract karama for the given user",
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
			fmt.Printf("%#v\n", e)

			return nil
		},
	}
}
