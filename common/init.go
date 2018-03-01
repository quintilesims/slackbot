package common

import "github.com/quintilesims/slackbot/db"

func Init(store db.Store) error {
	return initKarmaStore(store)
}
