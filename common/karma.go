package common

import "github.com/quintilesims/slackbot/db"

const StoreKeyKarma = "karma"

func initKarmaStore(store db.Store) error {
	karma := map[string]int{}
	if err := store.Read(StoreKeyKarma, &karma); err != nil {
		if _, ok := err.(db.MissingEntryError); ok {
			return store.Write(StoreKeyKarma, karma)
		}

		return err
	}

	return nil
}
