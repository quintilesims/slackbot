package runners

import (
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/lock"
	"github.com/quintilesims/slackbot/models"
	"github.com/quintilesims/slackbot/utils"
)

func NewRemindersRunner(l lock.Lock, store db.Store, client *slack.Client) *Runner {
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

		errs := []error{}
		now := time.Now().UTC()
		for reminderID, r := range reminders {
			if now.After(r.Time) {
				log.Printf("[DEBUG] [RemindersRunner] Sending reminder %s", reminderID)
				text := fmt.Sprintf("Hi %s! This is your reminder for %s", r.UserName, r.Message)
				if _, _, err := client.PostMessage(r.UserID, text, slack.PostMessageParameters{}); err != nil {
					log.Printf("[ERROR] [RemindersRunner] Failed to send reminder %s: %v", reminderID, err)
					errs = append(errs, err)
				}

				delete(reminders, reminderID)
			}
		}

		if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
			log.Printf("[ERROR] [RemindersRunner] Failed to save reminders: %v", err)
			errs = append(errs, err)
		}

		return utils.MultiError(errs)
	})
}
