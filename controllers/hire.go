package controllers

import (
	"strings"

	"github.com/nlopes/slack"
	"github.com/zpatrick/fireball"
)

type HireController struct{}

func NewHireController() *HireController {
	return &HireController{}
}

func (h *HireController) Routes() []*fireball.Route {
	routes := []*fireball.Route{
		{
			Path: "/hire",
			Handlers: fireball.Handlers{
				"POST": h.run,
			},
		},
	}

	return routes
}

// todo: validate token: https://github.com/nlopes/slack/blob/master/examples/slash/slash.go
func (h *HireController) run(c *fireball.Context) (fireball.Response, error) {
	s, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(s.Text) == "help" {
		m := slack.Msg{
			ResponseType: "ephemeral",
			Text:         "How to use /hire",
			Attachments: []slack.Attachment{
				{
					Text: "this is how it works",
				},
			},
		}

		return fireball.NewJSONResponse(200, m)
	}

	m := slack.Msg{
		ResponseType: "in_channel",
		Text:         "New Hire",
		Attachments: []slack.Attachment{
			{
				Text: "Select new hire",
				Actions: []slack.AttachmentAction{
					{
						Name:       "new_hire",
						Text:       "who is the nhire",
						Type:       "select",
						DataSource: "users",
					},
				},
			},
		},
	}

	return fireball.NewJSONResponse(200, m)
}
