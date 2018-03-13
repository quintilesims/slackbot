package bot

import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestKarmaTrackingBehavior(t *testing.T) {
	store := db.NewMemoryStore()
	karma := models.Karma{
		"dogs":  5,
		"cats":  -2,
		"sleep": 0,
	}

	if err := store.Write(models.StoreKeyKarma, karma); err != nil {
		t.Fatal(err)
	}

	events := []slack.RTMEvent{
		newSlackMessageEvent("dogs++"),
		newSlackMessageEvent("cats--"),
		newSlackMessageEvent("dogs++"),
		newSlackMessageEvent("sleep++"),
		newSlackMessageEvent("tacos++"),
		newSlackMessageEvent("sunday naps++"),
		newSlackMessageEvent("blah blah blah"),
		{},
	}

	b := NewKarmaTrackingBehavior(store)
	for _, e := range events {
		if err := b(e); err != nil {
			t.Fatal(err)
		}
	}

	result := models.Karma{}
	if err := store.Read(models.StoreKeyKarma, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Karma{
		"dogs":        7,
		"cats":        -3,
		"sleep":       1,
		"tacos":       1,
		"sunday naps": 1,
	}

	assert.Equal(t, expected, result)
}
