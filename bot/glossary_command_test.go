package bot

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestGlossaryAdd(t *testing.T) {
	store := newMemoryStore(t)

	cases := []struct {
		Input  string
		Output models.Glossary
	}{
		{
			Input:  "!glossary add foo one",
			Output: models.Glossary{"foo": "one"},
		},
		{
			Input:  "!glossary add bar two",
			Output: models.Glossary{"bar": "two", "foo": "one"},
		},
		{
			Input:  "!glossary add --force foo three",
			Output: models.Glossary{"bar": "two", "foo": "three"},
		},
	}

	for _, c := range cases {
		t.Run(c.Input, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewGlossaryCommand(store, w)
			if err := runTestApp(cmd, c.Input); err != nil {
				t.Fatal(err)
			}

			result := models.Glossary{}
			if err := store.Read(db.GlossaryKey, &result); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.Output, result)
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
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestGlossaryList(t *testing.T) {
	glossary := models.Glossary{
		"foo": "one",
		"bar": "two",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.GlossaryKey, glossary); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewGlossaryCommand(store, w)
	if err := runTestApp(cmd, "!glossary ls"); err != nil {
		t.Fatal(err)
	}

	for k, v := range glossary {
		assert.Contains(t, w.String(), k)
		assert.Contains(t, w.String(), v)
	}
}

func TestGlossaryListErrors(t *testing.T) {
	inputs := []string{
		"!glossary ls",
		"!glossary --count n ls",
	}

	cmd := NewGlossaryCommand(newMemoryStore(t), ioutil.Discard)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
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

func TestGlossaryShow(t *testing.T) {
	store := newMemoryStore(t)
	if err := store.Write(db.GlossaryKey, models.Glossary{"foo": "bar"}); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewGlossaryCommand(store, w)
	if err := runTestApp(cmd, "!glossary show foo"); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, w.String(), "foo")
	assert.Contains(t, w.String(), "bar")
}

func TestGlossaryShowErrors(t *testing.T) {
	inputs := []string{
		"!glossary show",
		"!glossary show foo",
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
