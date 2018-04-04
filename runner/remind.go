package runner

import (
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/quintilesims/slackbot/utils"
)

// InterviewReminders will be send out 1 hour before the interview
const InterviewReminderLead = time.Hour

// NewReminderRunner will return a runner that will send reminders to slack users.
// Each time the runner executes, it will read from the store and
func NewReminderRunner(store db.Store, client utils.SlackClient) *Runner {
	timers := []*time.Timer{}
	return &Runner{
		Name: "Remind",
		run: func() error {
			interviewTimers, err := getInterviewTimers(store, client)
			if err != nil {
				return err
			}

			// reset all of our timers
			for i := 0; i < len(timers); i++ {
				timers[i].Stop()
			}

			timers = interviewTimers
			return nil
		},
	}
}

func getInterviewTimers(store db.Store, client utils.SlackClient) ([]*time.Timer, error) {
	interviews := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &interviews); err != nil {
		return nil, err
	}

	timers := []*time.Timer{}
	for i := 0; i < len(interviews); i++ {
		d := time.Until(interviews[i].Time)
		if d < InterviewReminderLead {
			continue
		}

		timer := time.AfterFunc(d, func() {
			for _, interviewerID := range interviews[i].InterviewerIDs {
				text := fmt.Sprintf("Hi <@%s>! Just reminding you that ", interviewerID)
				text += fmt.Sprintf("you have an interview with *%s* ", interviews[i].Candidate)
				text += fmt.Sprintf(" at %s", interviews[i].Time.Format("03:04:05PM"))

				if _, _, _, err := client.SendMessage(interviewerID, slack.MsgOptionText(text, true)); err != nil {
					log.Printf("[ERROR] [Reminder] %v", err)
				}
			}
		})

		timers = append(timers, timer)
	}

	return timers, nil
}
