package slash

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

func NewInterviewCommand() *CommandSchema {
	return &CommandSchema{
		Name: "/interview",
		Help: "todo",
		Run: func(client *slack.Client, req slack.SlashCommand) (*slack.Msg, error) {
			args := strings.Split(req.Text, " ")
			if len(args) != 2 {
				return nil, NewSlackMessageError("Please specify the interviewee's name (one word) and date (dd/mm)")
			}

			msg := &slack.Msg{
				ResponseType: "in_channel",
				Text:         "New Interview",
				Attachments: []slack.Attachment{
					{
						Text:       "Who will be doing the interviews?",
						CallbackID: "123",
						Actions: []slack.AttachmentAction{
							{
								Name:       "Interview 1",
								Text:       "Who will be doing the first interview?",
								Type:       "select",
								DataSource: "users",
								//Confirm: true,
							},
							{
								Name:       "Interview 2",
								Text:       "Who will be doing the second interview?",
								Type:       "select",
								DataSource: "users",
							},
						},
					},
				},
			}

			return msg, nil
		},
		Callback: func(client *slack.Client, req slack.SlashCommand) (*slack.Msg, error) {
			return nil, fmt.Errorf("callback not impl!")
		},
	}
}
