package runner

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/mock"
	"github.com/quintilesims/slackbot/models"
)

func TestGetInterviewReminders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSlackClient := mock.NewMockSlackClient(ctrl)

	now := time.Now()
	interviews := models.Interviews{
		{Time: now.Add(InterviewReminderLead * 2).UTC(), InterviewerIDs: []string{"uid1", "uid2"}},
		{Time: now.Add(InterviewReminderLead * 2).UTC(), InterviewerIDs: []string{"uid3"}},
		{Time: now.Add(-InterviewReminderLead).UTC(), InterviewerIDs: []string{"bad"}},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.InterviewsKey, interviews); err != nil {
		t.Fatal(err)
	}

	c := make(chan bool)
	record := func(channel string, options ...slack.MsgOption) {
		c <- true
	}

	for _, id := range []string{"uid1", "uid2", "uid3"} {
		mockSlackClient.EXPECT().
			SendMessage(id, gomock.Any()).
			Do(record).
			Return("", "", "", nil)
	}

	timers, err := getInterviewTimers(store, mockSlackClient)
	if err != nil {
		t.Fatal(err)
	}

	// reset the timers to execute immediately
	for _, timer := range timers {
		timer.Reset(0)
	}

	// expect 3 calls: one for "uid1", "uid2", and "uid3"
	for i := 0; i < 3; i++ {
		select {
		case <-c:
		case <-time.After(time.Second):
			t.Fatalf("Timeout on index %d", i)
		}
	}
}
