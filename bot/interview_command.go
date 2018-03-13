package bot

import (
	"fmt"
	"io"
	"time"

	"github.com/quintilesims/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/urfave/cli"
)

// common time layouts
const (
	DateLayout = "01/02"
)

// NewInterviewCommand returns a cli.Command that manages !interview
func NewInterviewCommand(client slack.SlackClient, store db.Store, userID string, w io.Writer) cli.Command {
	return cli.Command{
		Name:  "!interview",
		Usage: "manage interviews",
		Subcommands: []cli.Command{
			{
				Name:      "add",
				Usage:     "add a new interview",
				ArgsUsage: "INTERVIEWEE DATE (mm/dd)",
				Action:    newInterviewAddAction(client, store, userID, w),
			},
		},
	}
}

func newInterviewAddAction(client slack.SlackClient, store db.Store, managerID string, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		interviewee := c.Args().Get(0)
		if interviewee == "" {
			return fmt.Errorf("Argument INTERVIEWEE is required")
		}

		date := c.Args().Get(1)
		if date == "" {
			return fmt.Errorf("Argument DATE is required")
		}

		t, err := time.Parse(DateLayout, date)
		if err != nil {
			return err
		}

		// normalize the time of the interview
		t = time.Date(0, t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		interview, err := addInterview(store, managerID, interviewee, t)
		if err != nil {
			return err
		}

		if err := addInterviewReminders(client, interview); err != nil {
			return err
		}

		text := fmt.Sprintf("Ok, I've added an interview for %s at %s", interviewee, t.Format(DateLayout))
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func addInterview(store db.Store, managerID, interviewee string, date time.Time) (models.Interview, error) {
	interview := models.Interview{
		ManagerID:   managerID,
		Interviewee: interviewee,
		Date:        date.UTC(),
	}

	interviews := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &interviews); err != nil {
		return models.Interview{}, err
	}

	interviews = append(interviews, interview)
	if err := store.Write(db.InterviewsKey, interviews); err != nil {
		return models.Interview{}, err
	}

	return interview, nil
}

func addInterviewReminders(client slack.SlackClient, i models.Interview) error {
	reminders := map[time.Time]string{
		i.Date.AddDate(0, 0, -1): fmt.Sprintf("Pre-interview items for %s", i.Interviewee),
		i.Date:                  fmt.Sprintf("Interview with %s today", i.Interviewee),
		i.Date.AddDate(0, 0, 1): fmt.Sprintf("Post-interview items for %s", i.Interviewee),
	}

	for t, text := range reminders {
		date := t.Format("01/02 at 9:00am")
		if err := client.AddReminder("", i.ManagerID, text, date); err != nil {
			return err
		}
	}

	return nil
}

func newInterviewListAction(client slack.SlackClient, store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return fmt.Errorf("Not implemented")
	}
}

func newInterviewRemoveAction(client slack.SlackClient, store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		return fmt.Errorf("Not implemented")
	}
}
