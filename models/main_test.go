package models

import (
	"fmt"

	"github.com/quintilesims/slack"
)

func newSlackMessage(user, format string, tokens ...interface{}) *slack.Message {
	return &slack.Message{
		Msg: slack.Msg{
			User: user,
			Text: fmt.Sprintf(format, tokens...),
		},
	}
}
