package bot

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/quintilesims/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/urfave/cli"
)

// NewInterviewCommand returns a cli.Command that manages !interview
func NewInterviewCommand(client slack.SlackClient, store db.Store, w io.Writer) cli.Command {
	return cli.Command{
		Name:  "!interview",
		Usage: "manage interviews",
		Subcommands: []cli.Command{
			{
				Name:      "add",
				Usage:     "add a new interview",
				ArgsUsage: "@MANAGER INTERVIEWEE",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "time",
						Value: "09:00AM",
						Usage: "time of the interview",
					},
					cli.StringFlag{
						Name:  "date",
						Value: time.Now().Format(DateLayout),
						Usage: "date of the interview",
					},
				},
				Action: newInterviewAddAction(client, store, w),
			},
			{
				Name:      "ls",
				Usage:     "list all interviews",
				ArgsUsage: " ",
				Action:    newInterviewListAction(store, w),
			},
			{
				Name:      "rm",
				Usage:     "remove an interview",
				ArgsUsage: "INTERVIEWEE DATE",
				Action:    newInterviewRemoveAction(store, w),
			},
		},
	}
}

func newInterviewAddAction(client slack.SlackClient, store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()
		escapedManager := args.Get(0)
		if escapedManager == "" {
			return fmt.Errorf("Argument @MANAGER is required")
		}

		interviewee := args.Get(1)
		if interviewee == "" {
			return fmt.Errorf("Argument INTERVIEWEE is required")
		}

		if len(args) > 2 {
			interviewee = strings.Join(args[1:], " ")
		}

		dateTimeInput := strings.ToUpper(c.String("date") + c.String("time"))
		date, err := time.ParseInLocation(DateTimeLayout, dateTimeInput, time.Local)
		if err != nil {
			return err
		}

		manager, err := parseSlackUser(client, escapedManager)
		if err != nil {
			return fmt.Errorf("Invalid argument for @MANAGER: %v", err)
		}

		interview := models.Interview{
			Manager: models.User{
				ID:   manager.ID,
				Name: manager.Name,
			},
			Interviewee: interviewee,
			Date:        date.UTC(),
		}

		if err := addInterview(store, interview); err != nil {
			return err
		}

		if err := addInterviewReminders(client, interview); err != nil {
			return err
		}

		text := fmt.Sprintf("Ok, I've added an interview for *%s* on %s\n",
			interviewee,
			date.Format(DateAtTimeLayout))
		text += fmt.Sprintf("I've also added a few reminders for <@%s> \n", manager.ID)
		text += "They can use `/remind list` to view their current reminders in Slack"

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func addInterview(store db.Store, interview models.Interview) error {
	interviews := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &interviews); err != nil {
		return err
	}

	interviews = append(interviews, interview)
	return store.Write(db.InterviewsKey, interviews)
}

func addInterviewReminders(client slack.SlackClient, interview models.Interview) error {
	date := interview.Date.Local()
	interviewee := interview.Interviewee
	reminders := map[time.Time]string{
		date.AddDate(0, 0, -1): fmt.Sprintf("Pre-interview items for %s", interviewee),
		date.AddDate(0, 0, 0):  fmt.Sprintf("Interview with %s", interviewee),
		date.AddDate(0, 0, 1):  fmt.Sprintf("Post-interview items for %s", interviewee),
	}

	for d, text := range reminders {
		dateStr := d.Format(DateAtTimeLayout)
		if err := client.AddReminder("", interview.Manager.ID, text, dateStr); err != nil {
			return err
		}
	}

	return nil
}

func newInterviewListAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		interviews := models.Interviews{}
		if err := store.Read(db.InterviewsKey, &interviews); err != nil {
			return err
		}

		if len(interviews) == 0 {
			return fmt.Errorf("I don't have any interviews scheduled")
		}

		text := "Here are the interviews I have:\n"
		for _, interview := range interviews {
			text += fmt.Sprintf("*%s* on %s (manager: <@%s>)\n",
				interview.Interviewee,
				interview.Date.Format(DateAtTimeLayout),
				interview.Manager.ID)
		}

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func newInterviewRemoveAction(store db.Store, w io.Writer) func(c *cli.Context) error {
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

		interviews := models.Interviews{}
		if err := store.Read(db.InterviewsKey, &interviews); err != nil {
			return err
		}

		var exists bool
		for i := 0; i < len(interviews); i++ {
			interview := interviews[i]
			if strings.ToLower(interview.Interviewee) == strings.ToLower(interviewee) &&
				interview.Date.Month() == t.Month() &&
				interview.Date.Day() == t.Day() {
				interviews = append(interviews[:i], interviews[i+1:]...)
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("No interviews with %s on %s found", interviewee, t.Format(DateLayout))
		}

		if err := store.Write(db.InterviewsKey, interviews); err != nil {
			return err
		}

		text := fmt.Sprintf("Ok, I've deleted the *%s* interview on %s", interviewee, t.Format(DateLayout))
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}
