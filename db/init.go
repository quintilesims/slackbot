package db

import "github.com/quintilesims/slackbot/models"

// Init will initialize the table entries for the specified store
func Init(store Store) error {
	entries := map[string]interface{}{
		models.StoreKeyCallbacks:  models.Callbacks{},
		models.StoreKeyChecklists: models.Checklists{},
		models.StoreKeyInterviews: models.Interviews{},
		models.StoreKeyLocks:      models.Locks{},
		models.StoreKeyReminders:  models.Reminders{},
	}

	for key := range entries {
		v := entries[key]
		if err := store.Read(key, &v); err != nil {
			if _, ok := err.(MissingEntryError); ok {
				if err := store.Write(key, v); err != nil {
					return err
				}

				continue
			}

			return err
		}
	}

	return nil
}
