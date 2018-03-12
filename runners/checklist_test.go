package runners

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestChecklistRunner(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
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

		assert.Equal(t, "uid", query.Get("channel"))
		assert.NotNil(t, query.Get("text"))

		resp := slack.SlackResponse{Ok: true}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatal(err)
		}

	}

	client, close := newSlackClient(handler)
	defer close()

	checklists := models.Checklists{
		"uid": models.Checklist{
			{IsChecked: true},
			{IsChecked: false},
		},
		"uid2": models.Checklist{
			{IsChecked: true},
		},
	}

	store := newMemoryStore(t)
	if err := store.Write(models.StoreKeyChecklists, checklists); err != nil {
		t.Fatal(err)
	}

	runner := NewChecklistRunner(store, client)
	if err := runner.Run(); err != nil {
		t.Fatal(err)
	}
}
