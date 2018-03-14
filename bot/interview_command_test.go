package bot

import (
	"bytes"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/mock"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestInterviewAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock.NewMockSlackClient(ctrl)
	client.EXPECT().
		GetUserInfo("manager_id").
		Return(&slack.User{ID: "manager_id", Name: "manager_name"}, nil)

	client.EXPECT().
		AddReminder("", "manager_id", gomock.Any(), gomock.Any()).
		Return(nil).
		Times(3)

	store := newMemoryStore(t)
	w := bytes.NewBuffer(nil)
	cmd := NewInterviewCommand(client, store, w)

	if err := runTestApp(cmd, "!interview add --date 12/31 --time 06:00am <@manager_id> John Doe"); err != nil {
		t.Fatal(err)
	}

	result := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Interviews{
		{
			Manager:     models.User{ID: "manager_id", Name: "manager_name"},
			Interviewee: "John Doe",
			Date:        time.Date(0, 12, 31, 6, 0, 0, 0, time.Local).UTC(),
		},
	}

	assert.Equal(t, expected, result)
}

func TestInterviewList(t *testing.T) {
	interviews := models.Interviews{
		{Interviewee: "John Doe"},
		{Interviewee: "Jane Doe"},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.InterviewsKey, interviews); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewInterviewCommand(nil, store, w)
	if err := runTestApp(cmd, "!interview ls"); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, w.String(), "John Doe")
	assert.Contains(t, w.String(), "Jane Doe")
}

func TestInterviewRemove(t *testing.T) {
	interviews := models.Interviews{
		{
			Interviewee: "John Doe",
			Date:        time.Date(0, 12, 31, 0, 0, 0, 0, time.Local).UTC(),
		},
		{
			Interviewee: "Jane Doe",
			Date:        time.Date(0, 12, 31, 0, 0, 0, 0, time.Local).UTC(),
		},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.InterviewsKey, interviews); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewInterviewCommand(nil, store, w)
	if err := runTestApp(cmd, "!interview rm \"John Doe\" 12/31"); err != nil {
		t.Fatal(err)
	}

	result := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Interviews{
		{
			Interviewee: "Jane Doe",
			Date:        time.Date(0, 12, 31, 0, 0, 0, 0, time.Local).UTC(),
		},
	}

	assert.Equal(t, expected, result)
}
