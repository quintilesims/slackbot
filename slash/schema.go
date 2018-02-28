package slash

import (
	"fmt"
	"strings"
)

type SlashCommandSchema struct {
	Command string
	// so slash commands have the first 'get'-like response, when the user enters like /command
	// then, it can handle up to 5 callbacks
	// a simple way to do this is to just have 2 functions for each schema:
	//
}

func (s *SlashCommandSchema) Validate() error {
	if !strings.HasPrefix(s.Command, "/") {
		return fmt.Errorf("Invalid command field '%s'. Slash commands must start with '/'", s.Command)
	}

	// todo: ensure both the 'get' function and 'callbvack' function is not nil
	return nil
}
