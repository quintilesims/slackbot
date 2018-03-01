package behaviors

import "github.com/nlopes/slack"

func newMessageRTMEvent(text string) slack.RTMEvent {
	return slack.RTMEvent{
		Data: &slack.MessageEvent{
			Msg: slack.Msg{Text: text},
		},
	}
}
