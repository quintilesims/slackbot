package rtm

import (
	"fmt"
	"io"

	"github.com/nlopes/slack"
)

func NewHelpBehavior(behaviors ...*BehaviorSchema) *BehaviorSchema {
	return &BehaviorSchema{
		Name:  "help",
		Usage: "help `command`",
		Help:  "Display help for the given command.",
		OnMessageEvent: func(e *slack.MessageEvent, w io.Writer) error {
			return fmt.Errorf("Not implemented")
		},
	}
}
