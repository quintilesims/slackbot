package bot

import (
	"io/ioutil"
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestCandidateAdd(t *testing.T) {
	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, models.Candidates{"John Doe": nil}); err != nil {
		t.Fatal(err)
	}

	cmd := NewCandidateCommand(store, ioutil.Discard)
	if err := runTestApp(cmd, "!candidate add --meta k1=v1 --meta k2=v2 Jane Doe"); err != nil {
		t.Fatal(err)
	}

	result := models.Candidates{}
	if err := store.Read(db.CandidatesKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Candidates{
		"John Doe": nil,
		"Jane Doe": map[string]string{
			"k1": "v1",
			"k2": "v2",
		},
	}

	assert.Equal(t, expected, result)
}

func TestCandidateAddErrors(t *testing.T) {
	inputs := []string{
		"!candidate add",
		"!candidate add John Doe",
		"!candidate add --meta stuff NAME",
		"!candidate add --meta key:val NAME",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.CandidatesKey, models.Candidates{"John Doe": nil}); err != nil {
		t.Fatal(err)
	}

	cmd := NewCandidateCommand(store, ioutil.Discard)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
