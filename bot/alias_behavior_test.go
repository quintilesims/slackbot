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

func TestAliasBehaviorInvalidate(t *testing.T) {
	store := db.NewMemoryStore()
	if err := store.Write(db.AliasesKey, models.Aliases{"input": "one"}); err != nil {
		t.Fatal(err)
	}

	b := NewAliasBehavior(store)

	e1 := newSlackMessageEvent("input")
	if err := b.Behavior()(e1); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "one", e1.Data.(*slack.MessageEvent).Text)
	if err := store.Write(db.AliasesKey, models.Aliases{"input": "two"}); err != nil {
		t.Fatal(err)
	}

	b.Invalidate()
	e2 := newSlackMessageEvent("input")
	if err := b.Behavior()(e2); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "two", e2.Data.(*slack.MessageEvent).Text)
}
