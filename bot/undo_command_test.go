package bot

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/mock"
)

func TestUndo(t *testing.T) {
	cases := map[string]struct {
		ChannelID string
		Setup     func(appClient *mock.MockSlackClient, history *slack.History)
	}{
		"Public Channel": {
			ChannelID: "C123",
			Setup: func(appClient *mock.MockSlackClient, history *slack.History) {
				appClient.EXPECT().
					GetChannelHistory("C123", gomock.Any()).
					Return(history, nil)
			},
		},
		"Private Channel": {
			ChannelID: "G123",
			Setup: func(appClient *mock.MockSlackClient, history *slack.History) {
				appClient.EXPECT().
					GetGroupHistory("G123", gomock.Any()).
					Return(history, nil)
			},
		},
		"Direct Message": {
			ChannelID: "D123",
			Setup: func(appClient *mock.MockSlackClient, history *slack.History) {
				appClient.EXPECT().
					GetIMHistory("D123", gomock.Any()).
					Return(history, nil)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			appClient := mock.NewMockSlackClient(ctrl)
			botClient := mock.NewMockSlackClient(ctrl)

			history := &slack.History{
				Messages: []slack.Message{
					newSlackMessage("usr_id", "timestamp1", ""),
					newSlackMessage("bot_id", "timestamp2", ""),
					newSlackMessage("bot_id", "timestamp3", ""),
					newSlackMessage("usr_id", "timestamp4", ""),
				},
			}

			c.Setup(appClient, history)

			botClient.EXPECT().
				DeleteMessage(c.ChannelID, "timestamp2").
				Return("", "", nil)

			cmd := NewUndoCommand(appClient, botClient, c.ChannelID, "bot_id")
			if err := runTestApp(cmd, "!undo"); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUndoError(t *testing.T) {
	cases := map[string]string{
		"bad channel": "some_channel",
		"no matches":  "G123",
	}

	for name, channelID := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			appClient := mock.NewMockSlackClient(ctrl)
			botClient := mock.NewMockSlackClient(ctrl)

			appClient.EXPECT().
				GetGroupHistory(gomock.Any(), gomock.Any()).
				Return(&slack.History{}, nil).
				AnyTimes()

			cmd := NewUndoCommand(appClient, botClient, channelID, "")
			if err := runTestApp(cmd, "!undo"); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
