package commands

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestRemindersAdd(t *testing.T) {
	now := time.Now()
	cases := map[string]models.Reminder{
		"!reminders add <@u1> foo": models.Reminder{
			UserID:  "u1",
			Message: "foo",
			Time:    time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, time.Local).UTC(),
		},
		"!reminders add <@u1> foo bar baz": models.Reminder{
			UserID:  "u1",
			Message: "foo bar baz",
			Time:    time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, time.Local).UTC(),
		},
		"!reminders add --date 05/06 <@u1> foo": models.Reminder{
			UserID:  "u1",
			Message: "foo",
			Time:    time.Date(now.Year(), 5, 6, 9, 0, 0, 0, time.Local).UTC(),
		},
		"!reminders add --date 12/31 <@u1> foo": models.Reminder{
			UserID:  "u1",
			Message: "foo",
			Time:    time.Date(now.Year(), 12, 31, 9, 0, 0, 0, time.Local).UTC(),
		},
		"!reminders add --time 01:23AM <@u1> foo": models.Reminder{
			UserID:  "u1",
			Message: "foo",
			Time:    time.Date(now.Year(), now.Month(), now.Day()+1, 1, 23, 0, 0, time.Local).UTC(),
		},
		"!reminders add --time 12:34PM <@u1> foo": models.Reminder{
			UserID:  "u1",
			Message: "foo",
			Time:    time.Date(now.Year(), now.Month(), now.Day()+1, 12, 34, 0, 0, time.Local).UTC(),
		},
		"!reminders add --time 01:23PM <@u1> foo": models.Reminder{
			UserID:  "u1",
			Message: "foo",
			Time:    time.Date(now.Year(), now.Month(), now.Day()+1, 13, 23, 0, 0, time.Local).UTC(),
		},
		"!reminders add --date 05/06 --time 01:23AM <@u1> foo": models.Reminder{
			UserID:  "u1",
			Message: "foo",
			Time:    time.Date(now.Year(), 5, 6, 1, 23, 0, 0, time.Local).UTC(),
		},
	}

	newID := func() string { return "r1" }
	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			store := newMemoryStore(t)
			cmd := NewRemindersCommand(store, ioutil.Discard, newID)
			if err := runTestApp(cmd, input); err != nil {
				t.Fatal(err)
			}

			reminders := models.Reminders{}
			if err := store.Read(models.StoreKeyReminders, &reminders); err != nil {
				t.Fatal(err)
			}

			// use string comparison since comparing time.Time objects is unreliable
			assert.Equal(t, expected.String(), reminders["r1"].String())
		})
	}
}

func TestRemindersAddErrors(t *testing.T) {
	inputs := []string{
		"!reminders add",
		"!reminders add <@user>",
		"!reminders add user",
		"!reminders add user message",
		"!reminders --date 1/23 add <@user> message",
		"!reminders --date 12/3 add <@user> message",
		"!reminders --date 01:23 add <@user> message",
		"!reminders --time 1pm add <@user> message",
		"!reminders --time 1:00pm add <@user> message",
		"!reminders --time 01:0pm add <@user> message",
		"!reminders --time 01/00pm add <@user> message",
		"!reminders --time 09:00 add <@user> message",
		"!reminders --time 14:00 add <@user> message",
	}

	store := newMemoryStore(t)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			cmd := NewRemindersCommand(store, ioutil.Discard, nil)
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
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

	store := newMemoryStore(t)
	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewRemindersCommand(store, w, nil)
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

	store := newMemoryStore(t)
	cmd := NewRemindersCommand(store, ioutil.Discard, nil)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestRemindersRemove(t *testing.T) {
	reminders := models.Reminders{
		"r1": models.Reminder{},
		"r2": models.Reminder{},
	}

	store := newMemoryStore(t)
	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		t.Fatal(err)
	}

	cmd := NewRemindersCommand(store, ioutil.Discard, nil)
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

	store := newMemoryStore(t)
	cmd := NewRemindersCommand(store, ioutil.Discard, nil)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
