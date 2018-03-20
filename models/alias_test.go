package models

import (
	"testing"

	"github.com/quintilesims/slack"
	"github.com/stretchr/testify/assert"
)

func TestAliasesApply(t *testing.T) {
	cases := map[string]struct {
		Aliases  Aliases
		Message  *slack.MessageEvent
		Expected string
	}{
		"Pass-through empty Alias": {
			Aliases:  Aliases{},
			Message:  newSlackMessageEvent("", "", "message"),
			Expected: "message",
		},
		"Pass-through no match": {
			Aliases:  Aliases{"message": ""},
			Message:  newSlackMessageEvent("", "", "other"),
			Expected: "other",
		},
		"Static Alias": {
			Aliases:  Aliases{"message": "Hello, World!"},
			Message:  newSlackMessageEvent("", "", "message"),
			Expected: "Hello, World!",
		},
		"Args-Templated Alias": {
			Aliases:  Aliases{"message": "{{ index .Args 0 }}"},
			Message:  newSlackMessageEvent("", "", "message arg0 arg1"),
			Expected: "arg0",
		},
		"ArgsString-Templated Alias": {
			Aliases:  Aliases{"message": "{{ .ArgsString }}"},
			Message:  newSlackMessageEvent("", "", "message arg0 arg1"),
			Expected: "arg0 arg1",
		},
		"Channel-Templated Alias": {
			Aliases:  Aliases{"message": "{{ .ChannelID }}"},
			Message:  newSlackMessageEvent("channel", "", "message"),
			Expected: "channel",
		},
		"User-Templated Alias": {
			Aliases:  Aliases{"message": "{{ .UserID }}"},
			Message:  newSlackMessageEvent("", "user", "message"),
			Expected: "user",
		},
		"Multi-Templated Alias": {
			Aliases:  Aliases{"message": "{{ .UserID }} in {{ .ChannelID }} says {{ .ArgsString }}"},
			Message:  newSlackMessageEvent("channel", "user", "message Hello, World!"),
			Expected: "user in channel says Hello, World!",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := c.Aliases.Apply(c.Message); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.Expected, c.Message.Text)
		})
	}
}
