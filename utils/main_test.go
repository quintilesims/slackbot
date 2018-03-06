package utils

import (
	"net/http"
	"net/http/httptest"

	"github.com/nlopes/slack"
)

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
