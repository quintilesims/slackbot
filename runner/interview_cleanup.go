package runner

import (
	"log"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
)

func NewInterviewCleanupRunner(store db.Store) *Runner {
	return &Runner{
		Name: "InterviewCleanup",
		run: func() error {
			interviews := models.Interviews{}
			if err := store.Read(db.InterviewsKey, &interviews); err != nil {
				return err
			}

			oneWeek := time.Hour * 24 * 7
			for i := 0; i < len(interviews); i++ {
				if time.Since(interviews[i].Date) >= oneWeek {
					log.Printf("[DEBUG] [InterviewCleanup] Removing old interview for %s", interviews[i].Interviewee)
					interviews = append(interviews[:i], interviews[i+1:]...)
					i--
				}
			}

			return store.Write(db.InterviewsKey, interviews)
		},
	}
}
