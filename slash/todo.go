package slash

import (
	"fmt"
	"strconv"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/models"
)

func NewTODOCommand() *CommandSchema {
	return &CommandSchema{
		Name: "/todo",
		Run: func(req slack.SlashCommand) (*slack.Message, error) {
			t := models.TODO{
				ID:   "tid",
				Text: "do thing",
			}

			// todo: load TODO list from store using s.UserID
			msg := &slack.Message{
				Msg: slack.Msg{
					Text: "top level text",
					Attachments: []slack.Attachment{
						newAttachmentForTODO(t),
					},
				},
			}

			return msg, nil
		},
		Callback: func(req slack.AttachmentActionCallback) (*slack.Message, error) {
			fmt.Printf("here")
			return nil, nil
		},
	}
}

func newAttachmentForTODO(t models.TODO) slack.Attachment {
	return slack.Attachment{
		Text:       t.Text,
		CallbackID: t.ID,
		Actions: []slack.AttachmentAction{
			{
				Name:  "check",
				Text:  "check",
				Type:  "button",
				Value: strconv.FormatBool(t.IsChecked),
			},
			{
				Name:  "delete",
				Text:  "delete",
				Type:  "button",
				Style: "danger",
				Value: strconv.FormatBool(t.IsChecked),
			},
		},
	}
}
