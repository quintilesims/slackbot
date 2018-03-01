package common

import "github.com/quintilesims/slackbot/db"

func Init(store db.Store) error {
	if err := initKarmaStore(store); err != nil {
		return err
	}

	if err := initRemindersStore(store); err != nil {
		return err
	}

	return nil
}
