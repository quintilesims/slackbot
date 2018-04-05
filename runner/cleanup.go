package runner

import (
	"log"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
)

// Interviews expire after one week
const InterviewExpiry = time.Hour * 24 * 7

// NewCleanupRunner returns a runner that removes old data from the specified store.
// This includes deleting interviews that are older than one week.
func NewCleanupRunner(store db.Store) *Runner {
	return &Runner{
		Name: "Cleanup",
		run: func() error {
			if err := cleanupInterviews(store); err != nil {
				return err
			}

			return nil
		},
	}
}

func cleanupInterviews(store db.Store) error {
	interviews := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &interviews); err != nil {
		return err
	}

	for i := 0; i < len(interviews); i++ {
		if time.Now().UTC().Sub(interviews[i].Time.UTC()) >= InterviewExpiry {
			log.Printf("[DEBUG] [Cleanup] Removing old interview %#v", interviews[i])
			interviews = append(interviews[:i], interviews[i+1:]...)
			i--
		}
	}

	return store.Write(db.InterviewsKey, interviews)
}
