package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/slash"
	"github.com/zpatrick/fireball"
)

type SlashCommandController struct {
	client   *slack.Client
	token    string
	commands []*slash.CommandSchema
}

func NewSlashCommandController(client *slack.Client, token string, commands ...*slash.CommandSchema) *SlashCommandController {
	return &SlashCommandController{
		client:   client,
		token:    token,
		commands: commands,
	}
}

func (s *SlashCommandController) Routes() []*fireball.Route {
	routes := []*fireball.Route{
		{
			Path: "/slash/run",
			Handlers: fireball.Handlers{
				"POST": s.run,
			},
		},
		{
			Path: "/slash/callback",
			Handlers: fireball.Handlers{
				"POST": s.callback,
			},
		},
	}

	return routes
}

// todo: validate token: https://github.com/nlopes/slack/blob/master/examples/slash/slash.go
func (s *SlashCommandController) run(c *fireball.Context) (fireball.Response, error) {
	req, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		return nil, err
	}

	for _, cmd := range s.commands {
		if cmd.Name == req.Command {
			args := strings.Split(req.Text, " ")
			if len(args) == 1 && args[0] == "help" {
				msg := slack.Msg{Text: cmd.Help}
				return fireball.NewJSONResponse(200, msg)
			}

			msg, err := cmd.Run(s.client, req)
			if err != nil {
				return nil, err
			}

			return fireball.NewJSONResponse(200, msg)
		}
	}

	return nil, slash.NewSlackMessageError("This bot is currently not setup to handle the '%s' command!", req.Command)
}

func (s *SlashCommandController) callback(c *fireball.Context) (fireball.Response, error) {
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
