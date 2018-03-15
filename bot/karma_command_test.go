package bot

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestKarma(t *testing.T) {
	karmas := models.Karmas{
		"dogs":   {Upvotes: 10, Downvotes: 0},
		"people": {Upvotes: 5, Downvotes: 5},
		"cats":   {Upvotes: 0, Downvotes: 10},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.KarmasKey, karmas); err != nil {
		t.Fatal(err)
	}

	cases := map[string][]string{
		"!karma *":                         []string{"dogs"},
		"!karma --ascending *":             []string{"cats"},
		"!karma --count 3 *":               []string{"dogs", "people", "cats"},
		"!karma --count 100 *":             []string{"dogs", "people", "cats"},
		"!karma --count 3 --ascending *":   []string{"cats", "people", "dogs"},
		"!karma --count 100 --ascending *": []string{"cats", "people", "dogs"},
		"!karma dogs":                      []string{"dogs"},
		"!karma *o*":                       []string{"dogs"},
		"!karma --ascending *o*":           []string{"people"},
		"!karma --count 2 *o*":             []string{"dogs", "people"},
		"!karma --count 100 *o*":           []string{"dogs", "people"},
	}

	for input, expectedMatches := range cases {
		t.Run(input, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewKarmaCommand(store, w)
			if err := runTestApp(cmd, input); err != nil {
				t.Fatal(err)
			}

			for _, expected := range expectedMatches {
				assert.Contains(t, w.String(), expected)
			}
		})
	}
}

func TestKarmaErrors(t *testing.T) {
	inputs := []string{
		"!karma",
		"!karma --ascending",
		"!karma --count 5",
		"!karma --count five *",
		"!karma keyThatDoesNotExist",
	}

	store := newMemoryStore(t)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			cmd := NewKarmaCommand(store, ioutil.Discard)
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
