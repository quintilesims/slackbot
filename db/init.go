package db

import "github.com/quintilesims/slackbot/models"

func Init(store Store) error {
	if err := initKarmaStore(store); err != nil {
		return err
	}

	if err := initLocksStore(store); err != nil {
		return err
	}

	if err := initRemindersStore(store); err != nil {
		return err
	}

	return nil
}

func initKarmaStore(store Store) error {
	karma := models.Karma{}
	if err := store.Read(models.StoreKeyKarma, &karma); err != nil {
		if _, ok := err.(MissingEntryError); ok {
			return store.Write(models.StoreKeyKarma, karma)
		}

		return err
	}

	return nil
}

func initLocksStore(store Store) error {
	locks := models.Locks{}
	if err := store.Read(models.StoreKeyLocks, &locks); err != nil {
		if _, ok := err.(MissingEntryError); ok {
			return store.Write(models.StoreKeyLocks, locks)
		}

		return err
	}

	return nil
}

func initRemindersStore(store Store) error {
	reminders := models.Reminders{}
	if err := store.Read(models.StoreKeyReminders, &reminders); err != nil {
		if _, ok := err.(MissingEntryError); ok {
			return store.Write(models.StoreKeyReminders, reminders)
		}

		return err
	}

	return nil
}
