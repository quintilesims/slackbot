package runner

import (
	"testing"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestInterviewCleanup(t *testing.T) {
	oneWeekAgo := -(time.Hour * 24 * 7)
	interviews := models.Interviews{
		{Interviewee: "John Doe", Date: time.Now().Add(oneWeekAgo - time.Minute)},
		{Interviewee: "Jane Doe", Date: time.Now().Add(oneWeekAgo + time.Minute)},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.InterviewsKey, interviews); err != nil {
		t.Fatal(err)
	}

	runner := NewInterviewCleanupRunner(store)
	if err := runner.Run(); err != nil {
		t.Fatal(err)
	}

	result := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &result); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result, 1)
	assert.Equal(t, "Jane Doe", result[0].Interviewee)
}
