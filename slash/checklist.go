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

func NewChecklistCommand(store db.Store) *CommandSchema {
	return &CommandSchema{
		Name:     "/checklist",
		Help:     "View/Manage your checklist with `/checklist`, or add an item with `/checklist add <message>`",
		Run:      newChecklistRun(store),
		Callback: newChecklistCallback(store),
	}
}

func newChecklistRun(store db.Store) func(slack.SlashCommand) (*slack.Message, error) {
	return func(req slack.SlashCommand) (*slack.Message, error) {
		args := strings.Split(req.Text, " ")
		switch {
		case len(args) == 0 || args[0] == "":
			return checklistShow(store, req.UserID)
		case args[0] == "add":
			return checklistAdd(store, req.UserID, strings.Join(args[1:], " "))
		default:
			return nil, NewSlackMessageError("Invalid usage: please use `/checklist help` for more information")
		}
	}
}

func newChecklistItemCallbackID(userID, itemID string) string {
	return fmt.Sprintf("%s.%s", userID, itemID)
}

func parseChecklistItemCallbackID(callbackID string) (string, string, error) {
	split := strings.SplitN(callbackID, ".", 2)
	if len(split) != 2 {
		return "", "", fmt.Errorf("Failed to parse checklist callback id: %s", callbackID)
	}

	return split[0], split[1], nil
}

func newAttachmentForChecklistItem(userID string, item models.ChecklistItem) slack.Attachment {
	checkAction := slack.AttachmentAction{
		Name: "check",
		Text: "check",
		Type: "button",
	}

	text := item.Text
	if item.IsChecked {
		text = fmt.Sprintf("~%s~", text)
		checkAction = slack.AttachmentAction{
			Name: "uncheck",
			Text: "uncheck",
			Type: "button",
		}
	}

	deleteAction := slack.AttachmentAction{
		Name:  "delete",
		Text:  "delete",
		Type:  "button",
		Style: "danger",
	}

	return slack.Attachment{
		Text:       text,
		Color:      "#3AA3E3",
		CallbackID: newChecklistItemCallbackID(userID, item.ID),
		Actions:    []slack.AttachmentAction{checkAction, deleteAction},
	}
}

func checklistShow(store db.Store, userID string) (*slack.Message, error) {
	checklists := models.Checklists{}
	if err := store.Read(models.StoreKeyChecklists, &checklists); err != nil {
		return nil, err
	}

	checklist, ok := checklists[userID]
	if !ok || len(checklist) == 0 {
		return nil, NewSlackMessageError("You currently don't have any items in your checklist")
	}

	attachments := make([]slack.Attachment, len(checklist))
	for i, item := range checklist {
		attachments[i] = newAttachmentForChecklistItem(userID, item)
	}

	msg := &slack.Message{
		Msg: slack.Msg{
			Text:        "*Your Checklist*",
			Attachments: attachments,
		},
	}

	return msg, nil
}

func checklistAdd(store db.Store, userID, message string) (*slack.Message, error) {
	if message == "" {
		return nil, NewSlackMessageError("MESSAGE is required")
	}

	checklists := models.Checklists{}
	if err := store.Read(models.StoreKeyChecklists, &checklists); err != nil {
		return nil, err
	}

	checklist, ok := checklists[userID]
	if !ok {
		checklist = models.Checklist{}
	}

	item := models.ChecklistItem{
		ID:        strconv.Itoa(int(time.Now().UnixNano())),
		Text:      message,
		IsChecked: false,
	}

	checklists[userID] = append(checklist, item)
	if err := store.Write(models.StoreKeyChecklists, checklists); err != nil {
		return nil, err
	}

	msg := &slack.Message{
		Msg: slack.Msg{
			Text: fmt.Sprintf("Ok! I've added '%s' to your checklist", message),
		},
	}

	return msg, nil
}

func newChecklistCallback(store db.Store) func(slack.AttachmentActionCallback) (*slack.Message, error) {
	newCheckFunc := func(isChecked bool, itemID string) func(checklist models.Checklist) (models.Checklist, error) {
		return func(checklist models.Checklist) (models.Checklist, error) {
			for i := range checklist {
				if checklist[i].ID == itemID {
					checklist[i].IsChecked = isChecked
					return checklist, nil
				}
			}

			return nil, NewSlackMessageError("That checklist item no longer exists!")
		}
	}

	newDeleteFunc := func(itemID string) func(checklist models.Checklist) (models.Checklist, error) {
		return func(checklist models.Checklist) (models.Checklist, error) {
			for i := range checklist {
				if checklist[i].ID == itemID {
					return append(checklist[:i], checklist[i+1:]...), nil
				}
			}

			return nil, NewSlackMessageError("That checklist item no longer exists!")
		}
	}

	return func(req slack.AttachmentActionCallback) (*slack.Message, error) {
		userID, itemID, err := parseChecklistItemCallbackID(req.CallbackID)
		if err != nil {
			return nil, err
		}

		funcs := []func(checklist models.Checklist) (models.Checklist, error){}
		for _, action := range req.Actions {
			switch action.Name {
			case "check":
				funcs = append(funcs, newCheckFunc(true, itemID))
			case "uncheck":
				funcs = append(funcs, newCheckFunc(false, itemID))
			case "delete":
				funcs = append(funcs, newDeleteFunc(itemID))
			default:
				return nil, fmt.Errorf("Unrecognized action name '%s'", action.Name)
			}
		}

		checklists := models.Checklists{}
		if err := store.Read(models.StoreKeyChecklists, &checklists); err != nil {
			return nil, err
		}

		checklist := checklists[userID]
		for _, fn := range funcs {
			checklist, err = fn(checklist)
			if err != nil {
				return nil, err
			}
		}

		checklists[userID] = checklist
		if err := store.Write(models.StoreKeyChecklists, checklists); err != nil {
			return nil, err
		}

		return checklistShow(store, req.User.ID)
	}
}
