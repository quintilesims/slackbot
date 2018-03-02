package commands

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestRemindersAdd(t *testing.T) {
}

func TestRemindersAddErrors(t *testing.T) {
}

func TestRemindersList(t *testing.T) {
	reminders := models.Reminders{
		"r1": models.Reminder{
			UserID:  "u1",
			Message: "message one",
			Time:    time.Date(0, 11, 5, 15, 45, 0, 0, time.UTC),
		},
		"r2": models.Reminder{
			UserID: "u2",
		},
	}

	store := db.NewMemoryStore()
	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewRemindersCommand(store, w)
	if err := runTestApp(cmd, "!reminders ls <@u1>"); err != nil {
		t.Fatal(err)
	}

	expected := "That user has the following reminders:\n"
	expected += "Reminder `r1`: message one at 03:45PM on 11/05\n"
	assert.Equal(t, expected, w.String())
}

func TestRemindersListErrors(t *testing.T) {
	inputs := []string{
		"!reminders ls",
		"!reminders ls user",
	}

	store := db.NewMemoryStore()
	if err := store.Write(models.StoreKeyReminders, models.Reminders{}); err != nil {
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

func TestRemindersRemove(t *testing.T) {
	reminders := models.Reminders{
		"r1": models.Reminder{},
		"r2": models.Reminder{},
	}

	store := db.NewMemoryStore()
	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		t.Fatal(err)
	}

	cmd := NewRemindersCommand(store, ioutil.Discard)
	if err := runTestApp(cmd, "!reminders rm r1"); err != nil {
		t.Fatal(err)
	}

	result := models.Reminders{}
	if err := store.Read(models.StoreKeyReminders, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Reminders{
		"r2": models.Reminder{},
	}

	assert.Equal(t, expected, result)
}

func TestRemindersRemoveErrors(t *testing.T) {
	inputs := []string{
		"!reminders rm",
		"!reminders rm r1",
	}

	store := db.NewMemoryStore()
	if err := store.Write(models.StoreKeyReminders, models.Reminders{}); err != nil {
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
