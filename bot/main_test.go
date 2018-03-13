package bot

import (
	"fmt"

	"github.com/nlopes/slack"
)

func newMessageRTMEvent(format string, tokens ...interface{}) slack.RTMEvent {
	return slack.RTMEvent{
		Data: &slack.MessageEvent{
			Msg: slack.Msg{Text: fmt.Sprintf(format, tokens...)},
		},
	}
}
