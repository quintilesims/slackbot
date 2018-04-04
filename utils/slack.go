package utils

import "github.com/nlopes/slack"

// SlackClient is used to mock the nlopes/slack client
type SlackClient interface {
	GetGroupHistory(group string, params slack.HistoryParameters) (*slack.History, error)
	GetChannelHistory(channelID string, params slack.HistoryParameters) (*slack.History, error)
	GetIMHistory(channel string, params slack.HistoryParameters) (*slack.History, error)
}
