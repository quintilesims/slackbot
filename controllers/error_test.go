package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quintilesims/slackbot/slash"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	cases := map[string]struct {
		Error          error
		ExpectedStatus int
	}{
		"standard error": {
			Error:          fmt.Errorf("some error"),
			ExpectedStatus: 500,
		},
		"slack message error": {
			Error:          slash.NewSlackMessageError("some error"),
			ExpectedStatus: 200,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ErrorHandler(recorder, &http.Request{}, c.Error)
			assert.Equal(t, c.ExpectedStatus, recorder.Code)
		})
	}
}
