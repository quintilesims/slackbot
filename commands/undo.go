package commands

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/urfave/cli"
)

// NewUndoCommand returns a cli.Command that manages !undo
func NewUndoCommand(client *slack.Client, channelID, botName string) cli.Command {
	return cli.Command{
		Name:  "!undo",
		Usage: "delete the last message sent by the slackbot",
		Action: func(c *cli.Context) error {
			// https://stackoverflow.com/questions/41111227/how-can-a-slack-bot-detect-a-direct-message-vs-a-message-in-a-channel

			fmt.Printf("Channel ID: %s\n", channelID)

			var getHistory func(string, slack.HistoryParameters) (*slack.History, error)
			switch {
			case strings.HasPrefix(channelID, "C"):
				getHistory = client.GetChannelHistory
			case strings.HasPrefix(channelID, "D"):
				getHistory = client.GetIMHistory
			case strings.HasPrefix(channelID, "G"):
				getHistory = client.GetGroupHistory
			default:
				return fmt.Errorf("Cannot find channel type for '%s'", channelID)
			}

			history, err := getHistory(channelID, slack.NewHistoryParameters())
			if err != nil {
				return err
			}

			var lastMessageTimestamp string
			for _, message := range history.Messages {
				if message.Username == strings.ToUpper(botName) {
					fmt.Printf("Message: %#v\n", message.Text)
					lastMessageTimestamp = message.Timestamp
					break
				}
			}

			fmt.Println("Bot name: ", botName)
			if lastMessageTimestamp == "" {
				return fmt.Errorf("Failed to find last message for %s", botName)
			}

			if _, _, err := client.DeleteMessage(channelID, lastMessageTimestamp); err != nil {
				return err
			}

			return nil
		},
	}
}
