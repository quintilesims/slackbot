package rtm

import (
	"fmt"
	"io"

	"github.com/nlopes/slack"
)

// todo: create pattern using flagSet (https://golang.org/pkg/flag/#FlagSet) to parse flags and args
// this may help with the quotation arg problem (e.g. !karma "string with spaces" vs.
func NewInterviewBehavior() *BehaviorSchema {
	return &BehaviorSchema{
		Name: "",
		Help: "",
		OnMessageEvent: func(e *slack.MessageEvent, w io.Writer) error {
			return fmt.Errorf("Interview Not Implemented")
		},
	}
}
