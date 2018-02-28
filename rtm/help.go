package rtm

import (
	"fmt"
	"io"
	"strings"

	"github.com/nlopes/slack"
)

func NewHelpBehavior(behaviors ...*BehaviorSchema) *BehaviorSchema {
	help := "Display help for a command using `!help {command}`"
	return &BehaviorSchema{
		Name: "help",
		Help: help,
		OnMessageEvent: func(e *slack.MessageEvent, w io.Writer) error {
			args := strings.Split(e.Msg.Text, " ")
			if len(args) == 0 || args[0] != "!help" {
				return nil
			}

			if len(args) == 1 || args[1] == "help" {
				if _, err := w.Write([]byte(help)); err != nil {
					return err
				}

				return nil
			}

			for _, b := range behaviors {
				if b.Name == args[1] {
					if _, err := w.Write([]byte(b.Help)); err != nil {
						return err
					}

					return nil
				}
			}

			return fmt.Errorf("Could not find any commands named '%s'", args[1])
		},
	}
}
