package runners

import (
	"log"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/lock"
	"github.com/quintilesims/slackbot/models"
)

func NewRemindersRunner(l lock.Lock, store db.Store) *Runner {
	return NewRunner("RemindersRunner", func() error {
		if err := l.Lock(false); err != nil {
			if _, ok := err.(lock.LockContentionError); ok {
				log.Printf("[INFO] [RemindersRunner] lock already acquired, stopping run")
			}

			return err
		}

		defer l.Unlock()

		reminders := models.Reminders{}
		if err := store.Read(models.StoreKeyReminders, &reminders); err != nil {
			return err
		}

		now := time.Now().UTC()
		for _, r := range reminders {
			if now.Before(r.Time) {
				panic("good!")
			}

		}

		return nil
	})
}
