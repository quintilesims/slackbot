package slash

import (
	"fmt"
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

func addInterviewChecklistItems(store db.Store, interviewID string, interview models.Interview) error {
	checklists := models.Checklists{}
	if err := store.Read(models.StoreKeyChecklists, &checklists); err != nil {
		return err
	}

	checklist, ok := checklists[interview.ManagerID]
	if !ok {
		checklist = models.Checklist{}
	}

	// todo: use guuiid generator
	// todo: have a time for reminders, or print reminders every day
	checklist = append(checklist,
		models.ChecklistItem{
			ID:     strconv.Itoa(int(time.Now().UnixNano())),
			Text:   fmt.Sprintf("Pre-interview stuff for %s", interview.Interviewee),
			Source: interviewID,
		},
		models.ChecklistItem{
			ID:     strconv.Itoa(int(time.Now().UnixNano())),
			Text:   fmt.Sprintf("Post-interview stuff for %s", interview.Interviewee),
			Source: interviewID,
		},
	)

	checklists[interview.ManagerID] = checklist
	return store.Write(models.StoreKeyChecklists, checklists)
}

func interviewAdd(store db.Store, args []string) (*slack.Message, error) {
	if len(args) < 3 {
		return nil, NewSlackMessageError("@MANAGER DATE and INTERVIEWEE are required")
	}

	managerID, managerName, err := parseEscapedUser(args[0])
	if err != nil {
		return nil, NewSlackMessageError("Invalid MANAGER: specify a manager by typing `@<username>`")
	}

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
	interview := models.Interview{
		ManagerID:   managerID,
		ManagerName: managerName,
		Interviewee: strings.Join(args[2:], " "),
		Date:        date,
	}

	interviews[interviewID] = interview
	if err := store.Write(models.StoreKeyInterviews, interviews); err != nil {
		return nil, err
	}

	if err := addInterviewChecklistItems(store, interviewID, interview); err != nil {
		return nil, err
	}

	msg := &slack.Message{
		Msg: slack.Msg{
			Text: fmt.Sprintf("Ok, I've added an interview for %s on %s", interview.Interviewee, date.Format(dateFormat)),
		},
	}

	return msg, nil
}

// this function is idempotent
func deleteInterviewChecklistItems(store db.Store, managerID, interviewID string) error {
	checklists := models.Checklists{}
	if err := store.Read(models.StoreKeyChecklists, &checklists); err != nil {
		return err
	}

	checklist, ok := checklists[managerID]
	if !ok {
		return nil
	}

	for i := 0; i < len(checklist); i++ {
		if checklist[i].Source == interviewID {
			checklist = append(checklist[:i], checklist[i+1:]...)
			i--
		}
	}

	checklists[managerID] = checklist
	return store.Write(models.StoreKeyChecklists, checklists)
}

func newInterviewCallback(store db.Store) func(slack.AttachmentActionCallback) (*slack.Message, error) {
	return func(req slack.AttachmentActionCallback) (*slack.Message, error) {
		interviewID := req.CallbackID
		// todo: delete checklist items

		interviews := models.Interviews{}
		if err := store.Read(models.StoreKeyInterviews, &interviews); err != nil {
			return nil, err
		}

		interview, ok := interviews[interviewID]
		if !ok {
			return nil, NewSlackMessageError("That interview no longer exists!")
		}

		if err := deleteInterviewChecklistItems(store, interview.ManagerID, interviewID); err != nil {
			return nil, err
		}

		delete(interviews, interviewID)
		if err := store.Write(models.StoreKeyInterviews, interviews); err != nil {
			return nil, err
		}

		return interviewsShow(store)
	}
}
