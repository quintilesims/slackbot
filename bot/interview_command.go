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
				ArgsUsage: "CANDIDATE DATE (mm/dd/yyyy) TIME (mm:hh{am/pm}) @INTERVIEWERS..",
				Action:    newInterviewAddAction(client, store, w),
			},
			{
				Name:      "ls",
				Usage:     "list interviews",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "count",
						Value: 10,
						Usage: "The maximum number of interviews to display",
					},
				},
				Action: newInterviewListAction(store, w),
			},
			{
				Name:      "rm",
				Usage:     "remove an interview",
				ArgsUsage: "CANDIDATE DATE (mm/dd/yyyy) TIME (mm:hh{am/pm})",
				Action:    newInterviewRemoveAction(store, w),
			},
		},
	}
}

func newInterviewAddAction(client slack.SlackClient, store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()
		candidate := args.Get(0)
		if candidate == "" {
			return fmt.Errorf("Argument CANDIDATE is required")
		}

		dateStr := args.Get(1)
		if dateStr == "" {
			return fmt.Errorf("Argument DATE is required")
		}

		timeStr := strings.ToUpper(args.Get(2))
		if timeStr == "" {
			return fmt.Errorf("Argument TIME is required")
		}

		interviewers := args[3:]
		if len(interviewers) == 0 {
			return fmt.Errorf("At least one interviewer is required")
		}

		t, err := time.Parse(DateTimeLayout, dateStr+" "+timeStr)
		if err != nil {
			return err
		}

		interviewerIDs := make([]string, len(interviewers))
		for i := 0; i < len(interviewers); i++ {
			interviewerID, err := parseEscapedUserID(interviewers[i])
			if err != nil {
				return err
			}

			interviewerIDs[i] = interviewerID
		}

		interviews := models.Interviews{}
		if err := store.Read(db.InterviewsKey, &interviews); err != nil {
			return err
		}

		interview := models.Interview{
			Candidate:      candidate,
			Time:           t,
			InterviewerIDs: interviewerIDs,
		}

		for i := 0; i < len(interviews); i++ {
			if interviews[i].Equals(interview) {
				return fmt.Errorf("An interview for that candidate on the same date and time already exists")
			}
		}

		interviews = append(interviews, interview)
		if err := store.Write(db.InterviewsKey, interviews); err != nil {
			return err
		}

		text := fmt.Sprintf("Ok, I've added an interview for *%s* on %s",
			interview.Candidate,
			interview.Time.Format(DateAtTimeLayout))

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
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
		for i := 0; i < len(interviews) && i < c.Int("count"); i++ {
			dateAtTime := interviews[i].Time.Format(DateAtTimeLayout)
			text += fmt.Sprintf("*%s* on %s with ", interviews[i].Candidate, dateAtTime)
			for _, interviewerID := range interviews[i].InterviewerIDs {
				text += fmt.Sprintf("<@%s> ", interviewerID)
			}

			text += "\n"
		}

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func newInterviewRemoveAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := c.Args()
		candidate := args.Get(0)
		if candidate == "" {
			return fmt.Errorf("Argument CANDIDATE is required")
		}

		dateStr := args.Get(1)
		if dateStr == "" {
			return fmt.Errorf("Argument DATE is required")
		}

		timeStr := strings.ToUpper(args.Get(2))
		if timeStr == "" {
			return fmt.Errorf("Argument TIME is required")
		}

		t, err := time.Parse(DateTimeLayout, dateStr+" "+timeStr)
		if err != nil {
			return err
		}

		interviews := models.Interviews{}
		if err := store.Read(db.InterviewsKey, &interviews); err != nil {
			return err
		}

		interview := models.Interview{
			Candidate: candidate,
			Time:      t,
		}

		var exists bool
		for i := 0; i < len(interviews); i++ {
			if interviews[i].Equals(interview) {
				exists = true
				interviews = append(interviews[:i], interviews[i+1:]...)
				i--
			}
		}

		if !exists {
			return fmt.Errorf("I couldn't find any interviews matching the specified name, date, and time")
		}

		if err := store.Write(db.InterviewsKey, interviews); err != nil {
			return err
		}

		text := fmt.Sprintf("Ok, I've deleted the interview for *%s*", candidate)
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}
