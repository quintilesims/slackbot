package slash

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

type CommandSchema struct {
	Name     string
	Help     string
	Run      func(client *slack.Client, req slack.SlashCommand) (*slack.Msg, error)
	Callback func(client *slack.Client, req slack.SlashCommand) (*slack.Msg, error)
}

func (s *CommandSchema) Validate() error {
	if !strings.HasPrefix(s.Name, "/") {
		return fmt.Errorf("Invalid Name '%s'. Names for slash commands must start with /", s.Name)
	}

	if s.Help == "" {
		return fmt.Errorf("Command '%s' cannot have an empty help field!", s.Help)
	}

	if s.Run == nil {
		return fmt.Errorf("Command '%s' cannot have a nil Run field!", s.Name)
	}

	if s.Callback == nil {
		return fmt.Errorf("Command '%s' cannot have a nil Callback field!", s.Name)
	}

	return nil
}
