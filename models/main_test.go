package models

import (
	"fmt"

	"github.com/quintilesims/slack"
)

func newSlackMessageEvent(user, format string, tokens ...interface{}) *slack.MessageEvent {
	return &slack.MessageEvent{
		Msg: slack.Msg{
			User: user,
			Text: fmt.Sprintf(format, tokens...),
		},
	}
}
