package rtm

import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/stretchr/testify/assert"
)

func TestKarmaInit(t *testing.T) {
	s := db.NewMemoryStore()
	if err := NewKarmaAction(s).Init(); err != nil {
		t.Fatal(err)
	}

	keys, err := s.Keys()
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, keys, StoreKeyKarma)
}

func TestKarmaUpdate(t *testing.T) {
	currentKarma := map[string]int{
		"cats":  -2,
		"sleep": 0,
		"tacos": 3,
		"dogs":  5,
	}

	s := db.NewMemoryStore()
	if err := s.Write(StoreKeyKarma, currentKarma); err != nil {
		t.Fatal(err)
	}

	events := []*slack.MessageEvent{
		newMessageEvent("dogs++"),
		newMessageEvent("cats--"),
		newMessageEvent("tacos++"),
		newMessageEvent("tacos++"),
		newMessageEvent("sleep--"),
		newMessageEvent("dogs++"),
		newMessageEvent("red++"),
		newMessageEvent("blue--"),
		newMessageEvent("sleep++"),
		newMessageEvent("cats++"),
		newMessageEvent("cats--"),
		newMessageEvent("sleep++"),
	}

	a := NewKarmaAction(s)
	for _, e := range events {
		if err := a.OnMessageEvent(e, nil); err != nil {
			t.Fatal(err)
		}
	}

	expected := map[string]int{
		"cats":  -3,
		"blue":  -1,
		"red":   1,
		"sleep": 1,
		"tacos": 5,
		"dogs":  7,
	}

	var result map[string]int
	if err := s.Read(StoreKeyKarma, &result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}
