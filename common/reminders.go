package common

import (
	"time"

	"github.com/quintilesims/slackbot/db"
)

const StoreKeyReminders = "reminders"

type Reminder struct {
	UserID  string
	Message string
	Time    time.Time
}

type Reminders map[string]Reminder

func initRemindersStore(store db.Store) error {
	reminders := Reminders{}
	if err := store.Read(StoreKeyReminders, &reminders); err != nil {
		if _, ok := err.(db.MissingEntryError); ok {
			return store.Write(StoreKeyReminders, reminders)
		}

		return err
	}

	return nil
}
