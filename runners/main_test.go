package runners

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func newMemoryStore(t *testing.T) *db.MemoryStore {
	store := db.NewMemoryStore()
	if err := db.Init(store); err != nil {
		t.Fatal(err)
	}

	return store
}

func newSlackClient(handler func(w http.ResponseWriter, r *http.Request)) (*slack.Client, func()) {
	server := httptest.NewServer(http.HandlerFunc(handler))
	slack.SLACK_API = server.URL + "/"
	client := slack.New("")
	close := func() {
		server.Close()
		slack.SLACK_API = "https://slack.com/api/"
	}

	return client, close
}
