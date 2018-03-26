package bot

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestGlossaryAdd(t *testing.T) {
	store := newMemoryStore(t)

	cases := map[string]models.Glossary{
		"!glossary add foo one":           models.Glossary{"foo": "one"},
		"!glossary add bar two":           models.Glossary{"bar": "two", "foo": "one"},
		"!glossary add --force foo three": models.Glossary{"bar": "two", "foo": "three"},
	}

	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewGlossaryCommand(store, w)
			if err := runTestApp(cmd, input); err != nil {
				t.Fatal(err)
			}

			result := models.Glossary{}
			if err := store.Read(db.GlossaryKey, &result); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, result)
		})
	}
}

func TestGlossaryAddErrors(t *testing.T) {
	glossary := models.Glossary{
		"foo": "one",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.GlossaryKey, glossary); err != nil {
		t.Fatal(err)
	}

	inputs := []string{
		"!glossary add",
		"!glossary add foo",
		"!glossary add foo one",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			cmd := NewGlossaryCommand(store, ioutil.Discard)
			if err := runTestApp(cmd, input); err == nil {
				fmt.Println(err)
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestGlossaryRemove(t *testing.T) {
	glossary := models.Glossary{
		"foo": "",
		"bar": "",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.GlossaryKey, glossary); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewGlossaryCommand(store, w)
	if err := runTestApp(cmd, "!glossary rm foo"); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, w.String(), "foo")

	result := models.Glossary{}
	if err := store.Read(db.GlossaryKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Glossary{
		"bar": "",
	}

	assert.Equal(t, expected, result)
}

func TestGlossaryRemoveErrors(t *testing.T) {
	inputs := []string{
		"!glossary rm",
		"!glossary rm foo",
	}

	store := newMemoryStore(t)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			cmd := NewGlossaryCommand(store, ioutil.Discard)
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestGlossarySearch(t *testing.T) {
	glossary := models.Glossary{
		"foo": "one",
		"bar": "two",
		"baz": "three",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.GlossaryKey, glossary); err != nil {
		t.Fatal(err)
	}

	cases := map[string]models.Glossary{
		"*":           {"foo": "one", "bar": "two"},
		"--count 1 *": {"bar": "two"},
		"b*":          {"bar": "two", "baz": "three"},
		"*a*":         {"bar": "two", "baz": "three"},
		"foo":         {"foo": "one"},
	}

	for glob, expected := range cases {
		t.Run(glob, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewGlossaryCommand(store, w)
			input := fmt.Sprintf("!glossary ls %s", glob)

			if err := runTestApp(cmd, input); err != nil {
				t.Fatal(err)
			}

			for key, val := range expected {
				assert.Contains(t, w.String(), key)
				assert.Contains(t, w.String(), val)
			}
		})
	}
}

func TestGlossarySearchErrors(t *testing.T) {
	inputs := []string{
		"!glossary search",
		"!glossary search foo",
	}

	store := newMemoryStore(t)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			cmd := NewGlossaryCommand(store, ioutil.Discard)
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
