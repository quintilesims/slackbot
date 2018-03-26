package models

import (
	"fmt"

	"github.com/quintilesims/slack"
)

func newSlackMessageEvent(channel, user, format string, tokens ...interface{}) *slack.MessageEvent {
	return &slack.MessageEvent{
		Msg: slack.Msg{
			Channel: channel,
			User:    user,
			Text:    fmt.Sprintf(format, tokens...),
		},
	}
}
