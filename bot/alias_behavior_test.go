package bot

import (
	"testing"

	"github.com/quintilesims/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestAliasBehavior(t *testing.T) {
	store := db.NewMemoryStore()
	transformers := models.Aliases{
		"x":   models.Alias{Pattern: "!x", Template: "!undo"},
		"say": models.Alias{Pattern: "!say *", Template: "{{ replace .Text \"!say\" \"!echo\"}}"},
	}

	if err := store.Write(db.AliasesKey, transformers); err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		Event  slack.RTMEvent
		Assert func(t *testing.T, e slack.RTMEvent)
	}{
		"Non-Message event": {
			Event:  slack.RTMEvent{},
			Assert: func(t *testing.T, e slack.RTMEvent) {},
		},
		"No matching patterns": {
			Event: newSlackMessageEvent("blah blah blah"),
			Assert: func(t *testing.T, e slack.RTMEvent) {
				assert.Equal(t, "blah blah blah", e.Data.(*slack.MessageEvent).Text)
			},
		},
		"!x": {
			Event: newSlackMessageEvent("!x"),
			Assert: func(t *testing.T, e slack.RTMEvent) {
				assert.Equal(t, "!undo", e.Data.(*slack.MessageEvent).Text)
			},
		},
		"!say": {
			Event: newSlackMessageEvent("!say Hello, World!"),
			Assert: func(t *testing.T, e slack.RTMEvent) {
				assert.Equal(t, "!echo Hello, World!", e.Data.(*slack.MessageEvent).Text)
			},
		},
	}

	b := NewAliasBehavior(store).Behavior()
	for name := range cases {
		t.Run(name, func(t *testing.T) {
			c := cases[name]
			if err := b(c.Event); err != nil {
				t.Fatal(err)
			}

			c.Assert(t, c.Event)
		})
	}
}

// todo: test invalidate
