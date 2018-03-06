package utils

import (
	"fmt"
	"regexp"

	"github.com/nlopes/slack"
)

// SlackUserParser will take an escaped slack user ID and convert it to a *slack.User object
// escaped slack user ID's are in the format <@id>
type SlackUserParser func(escaped string) (*slack.User, error)

// NewSlackUserParser returns a SlackUserParser that uses the specified client to lookup *slack.User objects 
func NewSlackUserParser(client *slack.Client) SlackUserParser {
	return func(escaped string) (*slack.User, error) {
		// escaped user format: '<@ABC123>'
		r := regexp.MustCompile("\\<\\@[a-zA-Z0-9]+\\>")
		if !r.MatchString(escaped) {
			return nil, fmt.Errorf("Invalid user: please enter a valid user by typing `@<username>`")
		}

		userID := escaped[2 : len(escaped)-1]
		return client.GetUserInfo(userID)
	}
}

// NewStaticUserParser returns a SlackUserParser that always returns a user with the specified ID and name
func NewStaticUserParser(id, name string) SlackUserParser {
	return func(escaped string) (*slack.User, error) {
		return &slack.User{ID: id, Name: name}, nil
	}
}
