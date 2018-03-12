package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/models"
	"github.com/quintilesims/slackbot/slash"
	"github.com/stretchr/testify/assert"
	"github.com/zpatrick/fireball"
)

func TestSlashCommandControllerRun(t *testing.T) {
	var called bool
	cmd := &slash.CommandSchema{
		Name: "!test",
		Run: func(slack.SlashCommand) (*slack.Message, error) {
			called = true
			return newSlackMessageWithCallback("callback_id"), nil
		},
	}

	form := url.Values{}
	form.Set("command", "!test")
	req := newFormRequest(t, form)
	c := &fireball.Context{Request: req}

	store := newMemoryStore(t)
	controller := NewSlashCommandController(store, cmd)
	resp, err := controller.run(c)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)

	callbacks := models.Callbacks{}
	if err := store.Read(models.StoreKeyCallbacks, &callbacks); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "!test", callbacks["callback_id"])
}

func TestSlashCommandControllerRunHelp(t *testing.T) {
	cmd := &slash.CommandSchema{
		Name: "!test",
		Help: "some help",
	}

	form := url.Values{}
	form.Set("command", "!test")
	form.Set("text", "help")
	req := newFormRequest(t, form)
	c := &fireball.Context{Request: req}

	store := newMemoryStore(t)
	controller := NewSlashCommandController(store, cmd)
	resp, err := controller.run(c)
	if err != nil {
		t.Fatal(err)
	}

	var result slack.Message
	recorder := unmarshalBody(t, resp, &result)
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "some help", result.Text)
}

func TestSlashCommandControllerRunError(t *testing.T) {
	form := url.Values{}
	form.Set("command", "!test")
	req := newFormRequest(t, form)
	c := &fireball.Context{Request: req}

	controller := NewSlashCommandController(newMemoryStore(t))
	_, err := controller.run(c)
	if _, ok := err.(*slash.SlackMessageError); !ok {
		t.Fatalf("Error was not SlackMessageError: %#v", err)
	}
}

func TestSlashCommandControllerCallback(t *testing.T) {
	var called bool
	cmd := &slash.CommandSchema{
		Name: "!test",
		Callback: func(slack.AttachmentActionCallback) (*slack.Message, error) {
			called = true
			return &slack.Message{}, nil
		},
	}

	callbacks := models.Callbacks{
		"callback_id": "!test",
	}

	store := newMemoryStore(t)
	if err := store.Write(models.StoreKeyCallbacks, callbacks); err != nil {
		t.Fatal(err)
	}

	encoded, err := json.Marshal(slack.AttachmentActionCallback{CallbackID: "callback_id"})
	if err != nil {
		t.Fatal(err)
	}

	escaped := url.QueryEscape(string(encoded))
	body := fmt.Sprintf("payload=%s", escaped)
	req, err := http.NewRequest("POST", "https://test.com/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	c := &fireball.Context{Request: req}
	controller := NewSlashCommandController(store, cmd)
	resp, err := controller.callback(c)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
	recorder := unmarshalBody(t, resp, nil)
	assert.Equal(t, 200, recorder.Code)
}

func TestSlashCommandControllerCallbackError(t *testing.T) {
	encoded, err := json.Marshal(slack.AttachmentActionCallback{CallbackID: "callback_id"})
	if err != nil {
		t.Fatal(err)
	}

	escaped := url.QueryEscape(string(encoded))
	body := fmt.Sprintf("payload=%s", escaped)
	req, err := http.NewRequest("POST", "https://test.com/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	c := &fireball.Context{Request: req}
	controller := NewSlashCommandController(newMemoryStore(t))
	if _, err := controller.callback(c); err == nil {
		t.Fatalf("Error was nil!")
	}
}
