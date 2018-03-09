package slash

import (
	"github.com/nlopes/slack"
)

type CommandSchema struct {
	Name     string
	Help     string
	Run      func(slack.SlashCommand) (*slack.Message, error)
	Callback func(slack.AttachmentActionCallback) (*slack.Message, error)
}

// todo: Validate() func?
