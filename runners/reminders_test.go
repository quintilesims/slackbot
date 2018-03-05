package runners

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/lock"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestReminders(t *testing.T) {
	lock := lock.NewMemoryLock()
	store := newMemoryStore(t)
	client, close := newSlackClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/chat.postMessage", r.URL.String())

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		query, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "u1", query.Get("channel"))
		assert.Equal(t, "REMINDER: some message", query.Get("text"))

		resp := slack.SlackResponse{Ok: true}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatal(err)
		}
	})
	defer close()

	reminders := models.Reminders{
		"r1": {
			UserID:  "u1",
			Message: "some message",
			Time:    time.Now().UTC(),
		},
		"r2": {
			UserID:  "u2",
			Message: "some other message",
			Time:    time.Now().Add(time.Hour).UTC(),
		},
	}

	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		t.Fatal(err)
	}

	runner := NewRemindersRunner(lock, store, client)
	if err := runner.Run(); err != nil {
		t.Fatal(err)
	}

	result := models.Reminders{}
	if err := store.Read(models.StoreKeyReminders, &result); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result, 1)
	assert.Equal(t, reminders["r2"].String(), result["r2"].String())
}
