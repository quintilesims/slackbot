package slash

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

type CommandSchema struct {
	Name string
	Run  func(client *slack.Client, req slack.SlashCommand) (*slack.Msg, error)
	// so slash commands have the first 'get'-like response, when the user enters like /command
	// then, it can handle up to 5 callbacks
	// a simple way to do this is to just have 2 functions for each schema:
	//
}

func (s *CommandSchema) Validate() error {
	if !strings.HasPrefix(s.Name, "/") {
		return fmt.Errorf("Invalid Name '%s'. Names for slash commands must start with /", s.Name)
	}

	// todo: ensure both the 'get' function and 'callbvack' function is not nil
	return nil
}
