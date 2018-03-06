package utils

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestSlackUserParser(t *testing.T) {
	client, close := newSlackClient(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Ok   bool
			User slack.User
		}{
			Ok: true,
			User: slack.User{
				ID:   "uid",
				Name: "uname",
			},
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatal(err)
		}
	})
	defer close()

	parser := NewSlackUserParser(client)
	user, err := parser("<@uid>")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "uid", user.ID)
	assert.Equal(t, "uname", user.Name)
}

func TestSlackUserParserError(t *testing.T) {
	inputs := []string{
		"@ABC123",
		"username",
	}

	parser := NewSlackUserParser(nil)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if _, err := parser(input); err == nil {
				t.Fatal("Error was nil")
			}
		})
	}
}
