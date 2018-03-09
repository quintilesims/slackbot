package slash

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
)

const dateFormat = "01/02"

func NewInterviewCommand(store db.Store) *CommandSchema {
	return &CommandSchema{
		Name:     "/interview",
		Help:     "View/Manage interviews with `/interview`, or add an interview with `/interview add @MANAGER mm/dd INTERVIEWEE`",
		Run:      newInterviewRun(store),
		Callback: newInterviewCallback(store),
	}
}

func newInterviewRun(store db.Store) func(slack.SlashCommand) (*slack.Message, error) {
	return func(req slack.SlashCommand) (*slack.Message, error) {
		args := strings.Split(req.Text, " ")
		switch {
		case len(args) == 0 || args[0] == "":
			return interviewsShow(store)
		case args[0] == "add":
			return interviewAdd(store, args[1:])
		default:
			return nil, NewSlackMessageError("Invalid usage: please use `/interview help` for more information")
		}
	}
}

func newAttachmentForInterview(interviewID string, i models.Interview) slack.Attachment {
	return slack.Attachment{
		Text:       fmt.Sprintf("*%s* on %s (manager: %s)", i.Interviewee, i.Date.Format(dateFormat), i.ManagerName),
		Color:      "#3AA3E3",
		CallbackID: interviewID,
		Actions: []slack.AttachmentAction{
			{
				Name:  "delete",
				Text:  "delete",
				Type:  "button",
				Style: "danger",
			},
		},
	}
}

func interviewsShow(store db.Store) (*slack.Message, error) {
	interviews := models.Interviews{}
	if err := store.Read(models.StoreKeyInterviews, &interviews); err != nil {
		return nil, err
	}

	if len(interviews) == 0 {
		return nil, NewSlackMessageError("I currently don't have any interviews scheduled")
	}

	attachments := make([]slack.Attachment, 0, len(interviews))
	for interviewID, interview := range interviews {
		attachments = append(attachments, newAttachmentForInterview(interviewID, interview))
	}

	msg := &slack.Message{
		Msg: slack.Msg{
			Text:        "Here are the interviews I currently have scheduled:",
			Attachments: attachments,
		},
	}

	return msg, nil
}

func interviewAdd(store db.Store, args []string) (*slack.Message, error) {
	if len(args) < 3 {
		return nil, NewSlackMessageError("@MANAGER DATE and INTERVIEWEE are required")
	}

	// escaped format:  <@U1234|user>
	escapedManager := args[0]
	r := regexp.MustCompile("\\<\\@[a-zA-Z0-9]+|[a-zA-Z0-9]+\\>")
	if !r.MatchString(escapedManager) {
		return nil, NewSlackMessageError("Invalid MANAGER: specify a manager by typing `@<username>`")
	}

	// strip '<@' from the front and '>' from the end
	split := strings.SplitN(escapedManager[2:len(escapedManager)-1], "|", 2)
	managerID := split[0]
	managerName := split[1]

	date, err := time.Parse(dateFormat, args[1])
	if err != nil {
		return nil, NewSlackMessageError("Invalid DATE: %v", err)
	}

	interviews := models.Interviews{}
	if err := store.Read(models.StoreKeyInterviews, &interviews); err != nil {
		return nil, err
	}

	// todo: use guid generator
	interviewID := strconv.Itoa(int(time.Now().UnixNano()))
	interviewee := strings.Join(args[2:], " ")
	interviews[interviewID] = models.Interview{
		ManagerID:   managerID,
		ManagerName: managerName,
		Interviewee: interviewee,
		Date:        date,
	}

	if err := store.Write(models.StoreKeyInterviews, interviews); err != nil {
		return nil, err
	}

	msg := &slack.Message{
		Msg: slack.Msg{
			Text: fmt.Sprintf("Ok, I've added an interview for %s on %s", interviewee, date.Format(dateFormat)),
		},
	}

	return msg, nil
}

func newInterviewCallback(store db.Store) func(slack.AttachmentActionCallback) (*slack.Message, error) {
	return func(req slack.AttachmentActionCallback) (*slack.Message, error) {
		interviewID := req.CallbackID
		// todo: delete checklist items

		interviews := models.Interviews{}
		if err := store.Read(models.StoreKeyInterviews, &interviews); err != nil {
			return nil, err
		}

		if _, ok := interviews[interviewID]; !ok {
			return nil, NewSlackMessageError("That interview no longer exists!")
		}

		delete(interviews, interviewID)
		if err := store.Write(models.StoreKeyInterviews, interviews); err != nil {
                        return nil, err
                }

		return interviewsShow(store)
	}
}
