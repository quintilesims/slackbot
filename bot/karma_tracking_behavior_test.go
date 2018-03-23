package bot

import (
	"testing"

	"github.com/quintilesims/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestKarmaTrackingBehavior(t *testing.T) {
	store := db.NewMemoryStore()
	karmas := models.Karmas{
		"dogs": models.Karma{Upvotes: 10, Downvotes: 0},
		"cats": models.Karma{Upvotes: 0, Downvotes: 10},
	}

	if err := store.Write(db.KarmasKey, karmas); err != nil {
		t.Fatal(err)
	}

	events := []slack.RTMEvent{
		newSlackMessageEvent("dogs++"),
		newSlackMessageEvent("dogs++"),
		newSlackMessageEvent("cats--"),
		newSlackMessageEvent("cats--"),
		newSlackMessageEvent("new++"),
		newSlackMessageEvent("new--"),
		newSlackMessageEvent("new+-"),
		newSlackMessageEvent("new-+"),
		newSlackMessageEvent("blah blah"),
		{},
	}

	b := NewKarmaTrackingBehavior(store)
	for _, e := range events {
		if err := b(e); err != nil {
			t.Fatal(err)
		}
	}

	result := models.Karmas{}
	if err := store.Read(db.KarmasKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Karmas{
		"dogs": models.Karma{Upvotes: 12, Downvotes: 0},
		"cats": models.Karma{Upvotes: 0, Downvotes: 12},
		"new":  models.Karma{Upvotes: 3, Downvotes: 3},
	}

	assert.Equal(t, expected, result)
}
