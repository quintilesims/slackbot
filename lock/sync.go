package sync

import (
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/utils"
)

// todo: sync isn't a great name; it's essentially a lock that doesn't wait; it just
// doesn't execute if the lock is under contention
type Sync interface {
	Do(key string, fn func() error) error
}

// example
func DoReminders(sync Sync, store db.Store) error {
	type Reminder struct {
		ID            func() string
		ShouldTrigger func(time.Time) bool
		Run           func() error
	}

	var reminders []Reminder
	if err := store.Read("some key", &reminders); err != nil {
		return err
	}

	errs := []error{}
	for _, r := range reminders {
		if r.ShouldTrigger(time.Now()) {
			if err := sync.Do(r.ID(), r.Run); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utils.MultiError(errs)
}
