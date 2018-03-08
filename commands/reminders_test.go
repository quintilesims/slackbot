package commands

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/quintilesims/slackbot/models"
	"github.com/quintilesims/slackbot/utils"
	"github.com/stretchr/testify/assert"
)

func TestRemindersAdd(t *testing.T) {
	now := time.Now()
	cases := map[string]models.Reminder{
		"!reminders add <@uid> foo": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo",
			Time:     time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, time.Local).UTC(),
		},
		"!reminders add <@uid> foo bar baz": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo bar baz",
			Time:     time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, time.Local).UTC(),
		},
		"!reminders add --date 05/06 <@uid> foo": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo",
			Time:     time.Date(now.Year(), 5, 6, 9, 0, 0, 0, time.Local).UTC(),
		},
		"!reminders add --date 12/31 <@uid> foo": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo",
			Time:     time.Date(now.Year(), 12, 31, 9, 0, 0, 0, time.Local).UTC(),
		},
		"!reminders add --time 01:23AM <@uid> foo": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo",
			Time:     time.Date(now.Year(), now.Month(), now.Day()+1, 1, 23, 0, 0, time.Local).UTC(),
		},
		"!reminders add --time 12:34PM <@uid> foo": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo",
			Time:     time.Date(now.Year(), now.Month(), now.Day()+1, 12, 34, 0, 0, time.Local).UTC(),
		},
		"!reminders add --time 01:23PM <@uid> foo": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo",
			Time:     time.Date(now.Year(), now.Month(), now.Day()+1, 13, 23, 0, 0, time.Local).UTC(),
		},
		"!reminders add --year 2050 <@uid> foo": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo",
			Time:     time.Date(2050, now.Month(), now.Day()+1, 9, 0, 0, 0, time.Local).UTC(),
		},
		"!reminders add --date 05/06 --time 01:23AM <@uid> foo": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo",
			Time:     time.Date(now.Year(), 5, 6, 1, 23, 0, 0, time.Local).UTC(),
		},
		"!reminders add --year 2050 --date 05/06 --time 01:23AM <@uid> foo": {
			UserID:   "uid",
			UserName: "uname",
			Message:  "foo",
			Time:     time.Date(2050, 5, 6, 1, 23, 0, 0, time.Local).UTC(),
		},
	}

	generateID := utils.NewStaticIDGenerator("rid")
	userParser := utils.NewStaticUserParser("uid", "uname")
	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			store := newMemoryStore(t)
			cmd := NewRemindersCommand(store, ioutil.Discard, generateID, userParser)
			if err := runTestApp(cmd, input); err != nil {
				t.Fatal(err)
			}

			reminders := models.Reminders{}
			if err := store.Read(models.StoreKeyReminders, &reminders); err != nil {
				t.Fatal(err)
			}

			// use string comparison since comparing time.Time objects is unreliable
			assert.Equal(t, expected.String(), reminders["rid"].String())
		})
	}
}

func TestRemindersAddErrors(t *testing.T) {
	inputs := []string{
		"!reminders add",
		"!reminders add user",
		"!reminders add --date 1/23 <@user> message",
		"!reminders add --date 12/3 <@user> message",
		"!reminders add --date 01:23 <@user> message",
		"!reminders add --time 1pm <@user> message",
		"!reminders add --time 1:00pm <@user> message",
		"!reminders add --time 01:0pm <@user> message",
		"!reminders add --time 01/00pm <@user> message",
		"!reminders add --time 09:00 <@user> message",
		"!reminders add --time 14:00 <@user> message",
		"!reminders add --year 50 <@user> message",
		"!reminders add --year 2010 <@user> message",
	}

	store := newMemoryStore(t)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			cmd := NewRemindersCommand(store, ioutil.Discard, nil, nil)
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestRemindersList(t *testing.T) {
	reminders := models.Reminders{
		"rid": models.Reminder{
			UserID:   "uid",
			UserName: "uname",
			Message:  "some message",
			Time:     time.Date(2010, 11, 5, 15, 45, 0, 0, time.UTC),
		},
		"r2": models.Reminder{
			UserID: "uid2",
		},
	}

	store := newMemoryStore(t)
	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	userParser := utils.NewStaticUserParser("uid", "uname")
	cmd := NewRemindersCommand(store, w, nil, userParser)
	if err := runTestApp(cmd, "!reminders ls <@uid>"); err != nil {
		t.Fatal(err)
	}

	expected := "uname has the following reminders:\n"
	expected += "Reminder `rid`: some message at 03:45PM on 11/05 2010\n"
	assert.Equal(t, expected, w.String())
}

func TestRemindersListErrors(t *testing.T) {
	inputs := []string{
		"!reminders ls",
	}

	store := newMemoryStore(t)
	cmd := NewRemindersCommand(store, ioutil.Discard, nil, nil)
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
		"rid1": models.Reminder{},
		"rid2": models.Reminder{},
	}

	store := newMemoryStore(t)
	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		t.Fatal(err)
	}

	cmd := NewRemindersCommand(store, ioutil.Discard, nil, nil)
	if err := runTestApp(cmd, "!reminders rm rid1"); err != nil {
		t.Fatal(err)
	}

	result := models.Reminders{}
	if err := store.Read(models.StoreKeyReminders, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Reminders{
		"rid2": models.Reminder{},
	}

	assert.Equal(t, expected, result)
}

func TestRemindersRemoveErrors(t *testing.T) {
	inputs := []string{
		"!reminders rm",
		"!reminders rm rid1",
	}

	store := newMemoryStore(t)
	cmd := NewRemindersCommand(store, ioutil.Discard, nil, nil)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
