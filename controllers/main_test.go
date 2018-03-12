package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/zpatrick/fireball"
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

func newSlackMessageWithCallback(callbackID string) *slack.Message {
	return &slack.Message{
		Msg: slack.Msg{
			Attachments: []slack.Attachment{
				{CallbackID: "callback_id"},
			},
		},
	}
}

func newFormRequest(t *testing.T, form url.Values) *http.Request {
	req, err := http.NewRequest("POST", "https://test.com/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	return req
}

func unmarshalBody(t *testing.T, resp fireball.Response, v interface{}) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	resp.Write(recorder, nil)

	if v != nil {
		if err := json.Unmarshal(recorder.Body.Bytes(), v); err != nil {
			t.Fatal(err)
		}
	}

	return recorder
}
