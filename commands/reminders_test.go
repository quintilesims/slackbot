package commands

import (
	"io/ioutil"
	"testing"

	"github.com/quintilesims/slackbot/common"
	"github.com/quintilesims/slackbot/db"
	"github.com/stretchr/testify/assert"
)

func TestRemindersAdd(t *testing.T) {
}

func TestRemindersAddErrors(t *testing.T) {
}

func TestRemindersList(t *testing.T) {
}

func TestRemindersListErrors(t *testing.T) {
}

func TestRemindersRemove(t *testing.T) {
	reminders := common.Reminders{
		"r1": common.Reminder{},
		"r2": common.Reminder{},
	}

	store := db.NewMemoryStore()
	if err := store.Write(common.StoreKeyReminders, reminders); err != nil {
		t.Fatal(err)
	}

	cmd := NewRemindersCommand(store, ioutil.Discard)
	if err := runTestApp(cmd, "!reminders rm r1"); err != nil {
		t.Fatal(err)
	}

	result := common.Reminders{}
	if err := store.Read(common.StoreKeyReminders, &result); err != nil {
		t.Fatal(err)
	}

	expected := common.Reminders{
		"r2": common.Reminder{},
	}

	assert.Equal(t, expected, result)
}

func TestRemindersRemoveErrors(t *testing.T) {
	inputs := []string{
		"!reminders rm",
		"!reminders rm r1",
	}

	store := db.NewMemoryStore()
	if err := store.Write(common.StoreKeyReminders, common.Reminders{}); err != nil {
		t.Fatal(err)
	}

	cmd := NewRemindersCommand(store, ioutil.Discard)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatalf("Error was nil!")
			}
		})
	}
}
