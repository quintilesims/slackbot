package utils

import (
	"fmt"
	"regexp"

	"github.com/nlopes/slack"
)

type SlackUserParser func(escaped string) (*slack.User, error)

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

func NewStaticUserParser(id, name string) SlackUserParser {
	return func(escaped string) (*slack.User, error) {
		return &slack.User{ID: id, Name: name}, nil
	}
}
