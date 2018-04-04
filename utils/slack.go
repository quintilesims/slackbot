package utils

import "github.com/nlopes/slack"

// SlackClient is used to mock the nlopes/slack client
type SlackClient interface {
	DeleteMessage(channel, messageTimestamp string) (string, string, error)
	GetGroupHistory(group string, params slack.HistoryParameters) (*slack.History, error)
	GetChannelHistory(channelID string, params slack.HistoryParameters) (*slack.History, error)
	GetIMHistory(channel string, params slack.HistoryParameters) (*slack.History, error)
	SendMessage(channel string, options ...slack.MsgOption) (string, string, string, error)
}
