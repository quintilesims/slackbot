package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/quintilesims/slackbot/slash"
	"github.com/zpatrick/fireball"
)

type SlashCommandController struct {
	store    db.Store
	commands []*slash.CommandSchema
}

func NewSlashCommandController(store db.Store, commands ...*slash.CommandSchema) *SlashCommandController {
	return &SlashCommandController{
		store:    store,
		commands: commands,
	}
}

func (s *SlashCommandController) Routes() []*fireball.Route {
	routes := []*fireball.Route{
		{
			Path: "/",
			Handlers: fireball.Handlers{
				"POST": s.run,
			},
		},
		{
			Path: "/callback",
			Handlers: fireball.Handlers{
				"POST": s.callback,
			},
		},
	}

	return routes
}

func (s *SlashCommandController) run(c *fireball.Context) (fireball.Response, error) {
	req, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		return nil, err
	}

	var cmd *slash.CommandSchema
	for _, command := range s.commands {
		if command.Name == req.Command {
			cmd = command
			break
		}
	}

	if cmd == nil {
		return nil, slash.NewSlackMessageError("No matching handler found for '%s'", req.Command)
	}

	if args := strings.Split(req.Text, " "); len(args) == 1 && args[0] == "help" {
		msg := &slack.Message{
			Msg: slack.Msg{
				Text: cmd.Help,
			},
		}

		return fireball.NewJSONResponse(200, msg)
	}

	msg, err := cmd.Run(req)
	if err != nil {
		return nil, err
	}

	callbacks := models.Callbacks{}
	if err := s.store.Read(models.StoreKeyCallbacks, &callbacks); err != nil {
		return nil, err
	}

	// map callback ids to a slash command name
	for _, a := range msg.Attachments {
		if a.CallbackID != "" {
			callbacks[a.CallbackID] = cmd.Name
		}
	}

	if err := s.store.Write(models.StoreKeyCallbacks, callbacks); err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, msg)
}

func (s *SlashCommandController) callback(c *fireball.Context) (fireball.Response, error) {
	req, err := parseCallback(c.Request.Body)
	if err != nil {
		return nil, err
	}

	callbacks := models.Callbacks{}
	if err := s.store.Read(models.StoreKeyCallbacks, &callbacks); err != nil {
		return nil, err
	}

	commandName, ok := callbacks[req.CallbackID]
	if !ok {
		return nil, fmt.Errorf("Not matching callback entry found for '%s'", req.CallbackID)
	}

	var cmd *slash.CommandSchema
	for _, command := range s.commands {
		if command.Name == commandName {
			cmd = command
			break
		}
	}

	if cmd == nil {
		return nil, slash.NewSlackMessageError("No matching handler found for '%s'", commandName)
	}

	msg, err := cmd.Callback(*req)
	if err != nil {
		return nil, err
	}

	return fireball.NewJSONResponse(200, msg)
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
