package behaviors

import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/common"
	"github.com/quintilesims/slackbot/db"
	"github.com/stretchr/testify/assert"
)

func TestKarmaTrackingBehavior(t *testing.T) {
	store := db.NewMemoryStore()
	karma := map[string]int{
		"dogs":  5,
		"cats":  -2,
		"sleep": 0,
	}

	if err := store.Write(common.StoreKeyKarma, karma); err != nil {
		t.Fatal(err)
	}

	events := []slack.RTMEvent{
		newMessageRTMEvent("dogs++"),
		newMessageRTMEvent("cats--"),
		newMessageRTMEvent("dogs++"),
		newMessageRTMEvent("sleep++"),
		newMessageRTMEvent("tacos++"),
		slack.RTMEvent{},
	}

	b := NewKarmaTrackingBehavior(store)
	for _, e := range events {
		if err := b(e); err != nil {
			t.Fatal(err)
		}
	}

	result := map[string]int{}
	if err := store.Read(common.StoreKeyKarma, &result); err != nil {
		t.Fatal(err)
	}

	expected := map[string]int{
		"dogs":  7,
		"cats":  -3,
		"sleep": 1,
		"tacos": 1,
	}

	assert.Equal(t, expected, result)
}
