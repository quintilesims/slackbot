package models

import (
	"testing"

	"github.com/quintilesims/slack"
	"github.com/stretchr/testify/assert"
)

func TestAliasApply(t *testing.T) {
	cases := map[string]struct {
		Alias    Alias
		Message  *slack.MessageEvent
		Expected string
	}{
		"Pass-through": {
			Alias:    Alias{},
			Message:  newSlackMessageEvent("user", "input"),
			Expected: "input",
		},
		"Static Transformation": {
			Alias:    Alias{Pattern: "*", Template: "output"},
			Message:  newSlackMessageEvent("user", "input"),
			Expected: "output",
		},
		"User-Templated Transformation": {
			Alias:    Alias{Pattern: "*", Template: "{{ .User }}"},
			Message:  newSlackMessageEvent("user", "input"),
			Expected: "user",
		},
		"Message-Templated Transformation": {
			Alias:    Alias{Pattern: "*", Template: "{{ .Text }}"},
			Message:  newSlackMessageEvent("user", "input"),
			Expected: "input",
		},
		"Patterned Transformation": {
			Alias:    Alias{Pattern: "*see you later*", Template: "Goodbye {{ .User }}!"},
			Message:  newSlackMessageEvent("user", "goodnight, I will see you later"),
			Expected: "Goodbye user!",
		},
		"Patterned Transformation pass-through": {
			Alias:    Alias{Pattern: "*see you later*", Template: "Goodbye {{ .User }}!"},
			Message:  newSlackMessageEvent("user", "see y'all later"),
			Expected: "see y'all later",
		},
		"Patterned & Templated Transformation": {
			Alias:    Alias{Pattern: "!say *", Template: "{{ replace .Text \"!say\" \"!echo\"}}"},
			Message:  newSlackMessageEvent("user", "!say Hello, World!"),
			Expected: "!echo Hello, World!",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := c.Alias.Apply(c.Message); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.Expected, c.Message.Text)
		})
	}
}
