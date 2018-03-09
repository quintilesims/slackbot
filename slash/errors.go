package slash

import (
	"fmt"

	"github.com/nlopes/slack"
)

type SlackMessageError struct {
	*slack.Msg
}

func NewSlackMessageError(format string, tokens ...interface{}) *SlackMessageError {
	return &SlackMessageError{
		&slack.Msg{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf(format, tokens...),
		},
	}
}

func (s *SlackMessageError) Error() string {
	return s.Text
}
