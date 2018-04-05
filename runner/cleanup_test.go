package runner

import (
	"testing"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestCleanupInterviews(t *testing.T) {
	now := time.Now().UTC()
	interviews := models.Interviews{
		{Candidate: "old1", Time: now.Add(-InterviewExpiry).UTC()},
		{Candidate: "old2", Time: now.Add(-InterviewExpiry * 2).UTC()},
		{Candidate: "new1", Time: now.UTC()},
		{Candidate: "new2", Time: now.Add(InterviewExpiry).UTC()},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.InterviewsKey, interviews); err != nil {
		t.Fatal(err)
	}

	if err := cleanupInterviews(store); err != nil {
		t.Fatal(err)
	}

	result := models.Interviews{}
	if err := store.Read(db.InterviewsKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Interviews{
		{Candidate: "new1", Time: now.UTC()},
		{Candidate: "new2", Time: now.Add(InterviewExpiry).UTC()},
	}

	assert.Equal(t, expected, result)
}
