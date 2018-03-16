package bot

import (
	"fmt"
	"strings"

	"github.com/quintilesims/slack"
	"github.com/urfave/cli"
)

// NewUndoCommand returns a cli.Command that manages !undo.
// Due to Slack's permission model, only clients which have "bot" authentication are allowed to delete messages sent by the bot.
// However, only clients which have "app" authentication are allowed to get chat history for Channels, Groups, and IMs.
// Because of this, we must pass in a SlackClient authenticated for each.
func NewUndoCommand(appClient, botClient slack.SlackClient, channelID, botID string) cli.Command {
	return cli.Command{
		Name:  "!undo",
		Usage: "delete the last message sent by the slackbot",
		Action: func(c *cli.Context) error {
			var getHistory func(string, slack.HistoryParameters) (*slack.History, error)
			switch {
			case strings.HasPrefix(channelID, "C"):
				getHistory = appClient.GetChannelHistory
			case strings.HasPrefix(channelID, "D"):
				getHistory = appClient.GetIMHistory
			case strings.HasPrefix(channelID, "G"):
				getHistory = appClient.GetGroupHistory
			default:
				return fmt.Errorf("Cannot find channel type for '%s'", channelID)
			}

			history, err := getHistory(channelID, slack.NewHistoryParameters())
			if err != nil {
				return err
			}

			var lastMessageTimestamp string
			for _, message := range history.Messages {
				if message.User == botID {
					lastMessageTimestamp = message.Timestamp
					break
				}
			}

			if lastMessageTimestamp == "" {
				return fmt.Errorf("Failed to find last message sent by this bot")
			}

			if _, _, err := botClient.DeleteMessage(channelID, lastMessageTimestamp); err != nil {
				return err
			}

			return nil
		},
	}
}
