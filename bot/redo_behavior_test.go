package bot

import (
	"testing"
	"time"

	"github.com/quintilesims/slack"
	"github.com/stretchr/testify/assert"
)

func TestRedoBehaviorRecord(t *testing.T) {
	events := map[string]slack.RTMEvent{
		"message event": newSlackMessageEvent("test"),
		"empty event":   {},
	}

	b := NewRedoBehavior(nil)
	for name, event := range events {
		t.Run(name, func(t *testing.T) {
			if err := b.Record("cid", event); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRedoBehaviorRecordErrors(t *testing.T) {
	events := map[string]slack.RTMEvent{
		"!redo": newSlackMessageEvent("!redo"),
	}

	b := NewRedoBehavior(nil)
	for name, event := range events {
		t.Run(name, func(t *testing.T) {
			if err := b.Record("cid", event); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestRedoBehaviorTrigger(t *testing.T) {
	c := make(chan slack.RTMEvent)
	b := NewRedoBehavior(c)

	expected := newSlackMessageEvent("test")
	if err := b.Record("cid", expected); err != nil {
		t.Fatal(err)
	}

	if err := b.Trigger("cid"); err != nil {
		t.Fatal(err)
	}

	select {
	case result := <-c:
		assert.Equal(t, expected, result)
	case <-time.After(time.Second):
		t.Fatal("Timeout!")
	}
}

func TestRedoBehaviorTriggerError(t *testing.T) {
	b := NewRedoBehavior(nil)
	if err := b.Trigger("cid"); err == nil {
		t.Fatal("Error was nil!")
	}
}
