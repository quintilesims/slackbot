package rtm

import (
	"io"
	"strings"

	"github.com/nlopes/slack"
)

func NewEchoBehavior() *BehaviorSchema {
	return &BehaviorSchema{
		Name:  "echo",
		Usage: "echo `[args...]`",
		Help:  "Display the given message.",
		OnMessageEvent: func(e *slack.MessageEvent, w io.Writer) error {
			args := strings.Split(e.Msg.Text, " ")
			if len(args) == 0 || args[0] != "!echo" {
				return nil
			}

			text := strings.Join(args[1:], " ")
			if _, err := w.Write([]byte(text)); err != nil {
				return err
			}

			return nil
		},
	}
}
