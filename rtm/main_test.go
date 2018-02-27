package rtm

import "github.com/nlopes/slack"

func newMessageEvent(text string) *slack.MessageEvent {
	return &slack.MessageEvent{
		Msg: slack.Msg{Text: text},
	}
}
