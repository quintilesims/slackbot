package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/nlopes/slack"
	"github.com/zpatrick/fireball"
)

type HireController struct {
	token string
}

func NewHireController(token string) *HireController {
	return &HireController{
		token: token,
	}
}

func (h *HireController) Routes() []*fireball.Route {
	routes := []*fireball.Route{
		{
			Path: "/",
			Handlers: fireball.Handlers{
				"POST": h.callback,
			},
		},
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

	/* todo: get verification token
	if !s.ValidateToken(h.token) {
		return nil, fireball.NewError(401, fmt.Errorf("Invalid token"), nil)
	}
	*/

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
				Text:       "Select new hire",
				CallbackID: "123",
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

func (h *HireController) callback(c *fireball.Context) (fireball.Response, error) {
	callback, err := parseCallback(c.Request.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", callback)

	return fireball.NewJSONResponse(200, "Thanks!")
}

func parseCallback(body io.ReadCloser) (*slack.AttachmentActionCallback, error) {
	// slack does something odd here, where instead of sending just json in
	// the body, they send "payload=<json>" with the json url encoded
	defer body.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, err
	}

	s := strings.TrimPrefix(buf.String(), "payload=")
	decodedJSON, err := url.QueryUnescape(s)
	if err != nil {
		return nil, err
	}

	var callback *slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(decodedJSON), &callback); err != nil {
		return nil, err
	}

	return callback, nil
}
