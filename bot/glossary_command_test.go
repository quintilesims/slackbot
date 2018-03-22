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

func TestGlossaryDefine(t *testing.T) {
	inputs := map[string]string{
		"foo": "bar",
		"bar": "baz is multiple words",
		"baz": "foo",
	}

	store := newMemoryStore(t)
	for key, val := range inputs {
		w := bytes.NewBuffer(nil)
		cmd := NewGlossaryCommand(store, w)
		input := fmt.Sprintf("!glossary define %s %s", key, val)

		if err := runTestApp(cmd, input); err != nil {
			t.Fatal(err)
		}

		assert.Contains(t, w.String(), key)
		assert.Contains(t, w.String(), val)
	}

	result := models.Glossary{}
	if err := store.Read(db.GlossaryKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Glossary{
		"foo": "bar",
		"bar": "baz is multiple words",
		"baz": "foo",
	}

	assert.Equal(t, expected, result)
}

func TestGlossaryDefineErrors(t *testing.T) {
	inputs := []string{
		"!glossary define",
		"!glossary define foo",
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

func TestGlossaryRemove(t *testing.T) {
	glossary := models.Glossary{
		"foo": "",
		"bar": "",
		"baz": "",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.GlossaryKey, glossary); err != nil {
		t.Fatal(err)
	}

	for key := range glossary {
		w := bytes.NewBuffer(nil)
		cmd := NewGlossaryCommand(store, w)
		input := fmt.Sprintf("!glossary rm %s", key)

		if err := runTestApp(cmd, input); err != nil {
			t.Fatal(err)
		}

		assert.Contains(t, w.String(), key)
	}

	result := models.Glossary{}
	if err := store.Read(db.GlossaryKey, &result); err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result, 0)
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
		"foo": "bar",
		"bar": "baz",
		"baz": "foo",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.GlossaryKey, glossary); err != nil {
		t.Fatal(err)
	}

	cases := map[string]models.Glossary{
		"*":           {"foo": "bar", "bar": "baz", "baz": "foo"},
		"--count 2 *": {"foo": "bar", "bar": "baz", "baz": "foo"},
		"b*":          {"bar": "baz", "baz": "foo"},
		"*a*":         {"bar": "baz", "baz": "foo"},
		"foo":         {"foo": "bar"},
	}

	for glob, expected := range cases {
		t.Run(glob, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewGlossaryCommand(store, w)
			input := fmt.Sprintf("!glossary search %s", glob)

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
