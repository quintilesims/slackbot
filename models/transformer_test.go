package models

import (
	"testing"

	"github.com/quintilesims/slack"
	"github.com/stretchr/testify/assert"
)

func TestTransformerApply(t *testing.T) {
	cases := map[string]struct {
		Transformer Transformer
		Message     *slack.Message
		Expected    string
	}{
		"Pass-through": {
			Transformer: Transformer{},
			Message:     newSlackMessage("user", "input"),
			Expected:    "input",
		},
		"Static Transformation": {
			Transformer: Transformer{Pattern: "*", Template: "output"},
			Message:     newSlackMessage("user", "input"),
			Expected:    "output",
		},
		"User-Templated Transformation": {
			Transformer: Transformer{Pattern: "*", Template: "{{ .User }}"},
			Message:     newSlackMessage("user", "input"),
			Expected:    "user",
		},
		"Message-Templated Transformation": {
			Transformer: Transformer{Pattern: "*", Template: "{{ .Text }}"},
			Message:     newSlackMessage("user", "input"),
			Expected:    "input",
		},
		"Patterned Transformation": {
			Transformer: Transformer{Pattern: "*see you later*", Template: "Goodbye {{ .User }}!"},
			Message:     newSlackMessage("user", "goodnight, I will see you later"),
			Expected:    "Goodbye user!",
		},
		"Patterned Transformation pass-through": {
			Transformer: Transformer{Pattern: "*see you later*", Template: "Goodbye {{ .User }}!"},
			Message:     newSlackMessage("user", "see y'all later"),
			Expected:    "see y'all later",
		},
		"Patterned & Templated Transformation": {
			Transformer: Transformer{Pattern: "!say *", Template: "{{ replace .Text \"!say\" \"!echo\"}}"},
			Message:     newSlackMessage("user", "!say Hello, World!"),
			Expected:    "!echo Hello, World!",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := c.Transformer.Apply(c.Message); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.Expected, c.Message.Text)
		})
	}
}
