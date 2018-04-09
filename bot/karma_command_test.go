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

	cases := map[string]struct {
		Input  string
		Output []string
	}{
		"single count": {
			Input:  "!karma --count 1 *",
			Output: []string{"dogs"},
		},
		"single count ascending": {
			Input:  "!karma --count 1 --ascending *",
			Output: []string{"cats"},
		},
		"three count": {
			Input:  "!karma --count 3 *",
			Output: []string{"dogs", "people", "cats"},
		},
		"one-hundred count": {
			Input:  "!karma --count 100 *",
			Output: []string{"dogs", "people", "cats"},
		},
		"three count ascending": {
			Input:  "!karma --count 3 --ascending *",
			Output: []string{"dogs", "people", "cats"},
		},
		"one-hundred count ascending": {
			Input:  "!karma --count 100 --ascending *",
			Output: []string{"cats", "people", "dogs"},
		},
		"exact match": {
			Input:  "!karma dogs",
			Output: []string{"dogs"},
		},
		"wildcards": {
			Input:  "!karma *o*",
			Output: []string{"dogs"},
		},
		"wildcards ascending": {
			Input:  "!karma --ascending *o*",
			Output: []string{"people"},
		},
		"wildcards count": {
			Input:  "!karma --count 2 *o*",
			Output: []string{"dogs", "people"},
		},
		"wildcards count one-hundred": {
			Input:  "!karma --count 2 *o*",
			Output: []string{"dogs", "people"},
		},
	}

	for name := range cases {
		t.Run(name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewKarmaCommand(store, w)
			if err := runTestApp(cmd, cases[name].Input); err != nil {
				t.Fatal(err)
			}

			// TODO: Need to test orders
			for _, expected := range cases[name].Output {
				assert.Contains(t, w.String(), expected)
			}
		})
	}
}

func TestKarmaErrors(t *testing.T) {
	cases := map[string]string{
		"missing GLOB":           "!karma",
		"missing GLOB ascending": "!karma --ascending",
		"missing GlOB count":     "!karma --count 5",
		"missing int count":      "!karma --count five *",
		"unmatched key":          "!karma keyThatDoesNotExist",
	}

	store := newMemoryStore(t)
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			cmd := NewKarmaCommand(store, ioutil.Discard)
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
