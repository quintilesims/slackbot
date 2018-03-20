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
	aliases := models.Aliases{
		"input": "output",
	}

	if err := store.Write(db.AliasesKey, aliases); err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		Event  slack.RTMEvent
		Assert func(t *testing.T, e slack.RTMEvent)
	}{
		"Non-Message event": {
			Event: slack.RTMEvent{},
			Assert: func(t *testing.T, e slack.RTMEvent) {
				assert.Equal(t, slack.RTMEvent{}, e)
			},
		},
		"No matching alias": {
			Event: newSlackMessageEvent("blah blah blah"),
			Assert: func(t *testing.T, e slack.RTMEvent) {
				assert.Equal(t, "blah blah blah", e.Data.(*slack.MessageEvent).Text)
			},
		},
		"input to output": {
			Event: newSlackMessageEvent("input"),
			Assert: func(t *testing.T, e slack.RTMEvent) {
				assert.Equal(t, "output", e.Data.(*slack.MessageEvent).Text)
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
