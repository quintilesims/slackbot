package slash

import (
	"fmt"

	"github.com/nlopes/slack"
)

type SlackMessageError struct {
	*slack.Msg
}

func NewSlackMessageError(ephemeral bool, format string, tokens ...interface{}) *SlackMessageError {
	responseType := "ephemeral"
	if !ephemeral {
		responseType = "in_channel"
	}

	msg := &slack.Msg{
		ResponseType: responseType,
		Text:         fmt.Sprintf(format, tokens...),
	}

	msg.Text += "\nPlease contact the bot administrator"
	return &SlackMessageError{msg}
}

func (s *SlackMessageError) Error() string {
	return s.Text
}
