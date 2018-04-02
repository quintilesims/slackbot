package bot

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestInterviewAdd(t *testing.T) {
	store := newMemoryStore(t)
	w := bytes.NewBuffer(nil)
	cmd := NewInterviewCommand(store, w)

	if err := runTestApp(cmd, "!interview add \"John Doe\" 03/15/2014 09:35am <@uid1> <@uid2>"); err != nil {
		t.Fatal(err)
	}

	result := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Interviews{
		{
			Candidate:      "John Doe",
			InterviewerIDs: []string{"uid1", "uid2"},
			Time:           time.Date(2014, 3, 15, 9, 35, 0, 0, time.Local).UTC(),
		},
	}

	assert.Equal(t, expected, result)
}

func TestInterviewAddErrors(t *testing.T) {
	inputs := []string{
		"!interview add",
		"!interview add NAME",
		"!interview add NAME 03/15/2006",
		"!interview add NAME 03/15/2006 09:00am",
		"!interview add NAME 03/15/2006 09:00am uname",
		"!interview add NAME 03/15/2006 09:00am @uname",
		"!interview add NAME 15/03/2006 09:00am <@uid>",
		"!interview add NAME 3/15/2006 09:00am <@uid>",
		"!interview add NAME 03/15/06 09:00am <@uid>",
		"!interview add NAME 03/15/2006 9 <@uid>",
		"!interview add NAME 03/15/2006 9am <@uid>",
		"!interview add NAME 03/15/2006 9:00am <@uid>",
		"!interview add NAME 03/15/2006 09:00 <@uid>",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			cmd := NewInterviewCommand(newMemoryStore(t), ioutil.Discard)
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestInterviewList(t *testing.T) {
	interviews := models.Interviews{
		{Candidate: "John Doe"},
		{Candidate: "Jane Doe"},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.InterviewsKey, interviews); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewInterviewCommand(store, w)
	if err := runTestApp(cmd, "!interview ls"); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, w.String(), "John Doe")
	assert.Contains(t, w.String(), "Jane Doe")
}

func TestInterviewListError(t *testing.T) {
	cmd := NewInterviewCommand(newMemoryStore(t), ioutil.Discard)
	if err := runTestApp(cmd, "!interview ls"); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestInterviewRemove(t *testing.T) {
	interviews := models.Interviews{
		{
			Candidate: "John Doe",
			Time:      time.Date(2006, 12, 31, 9, 0, 0, 0, time.Local).UTC(),
		},
		{
			Candidate: "John Doe",
			Time:      time.Date(2006, 12, 31, 14, 0, 0, 0, time.Local).UTC(),
		},
		{
			Candidate: "Jane Doe",
			Time:      time.Date(2006, 12, 31, 9, 0, 0, 0, time.Local).UTC(),
		},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.InterviewsKey, interviews); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewInterviewCommand(store, w)
	if err := runTestApp(cmd, "!interview rm \"John Doe\" 12/31/2006 09:00am"); err != nil {
		t.Fatal(err)
	}

	result := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Interviews{
		{
			Candidate: "John Doe",
			Time:      time.Date(2006, 12, 31, 14, 0, 0, 0, time.Local).UTC(),
		},
		{
			Candidate: "Jane Doe",
			Time:      time.Date(2006, 12, 31, 9, 0, 0, 0, time.Local).UTC(),
		},
	}

	assert.ElementsMatch(t, expected, result)
}

func TestInterviewRemoveErrors(t *testing.T) {
	inputs := []string{
		"!interview rm",
		"!interview rm John",
		"!interview rm John 03/15/2006",
		"!interview rm John 03/15/2006 09:00am",
		"!interview rm John 15/03/2006 09:00am",
		"!interview rm John 3/15/2006 09:00am",
		"!interview rm John 03/15/06 09:00am",
		"!interview rm John 03/15/2006 9",
		"!interview rm John 03/15/2006 9am",
		"!interview rm John 03/15/2006 9:00am",
		"!interview rm John 03/15/2006 09:00",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			cmd := NewInterviewCommand(newMemoryStore(t), ioutil.Discard)
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
