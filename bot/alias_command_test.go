package bot

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestAliasAdd(t *testing.T) {
	var called bool
	invalidate := func() {
		called = true
	}

	store := newMemoryStore(t)
	w := bytes.NewBuffer(nil)
	cmd := NewAliasCommand(store, w, invalidate)
	if err := runTestApp(cmd, "!alias add !foo \"!echo Hello, World!\""); err != nil {
		t.Fatal(err)
	}

	result := models.Aliases{}
	if err := store.Read(db.AliasesKey, &result); err != nil {
		t.Fatal(err)
	}

	expected := models.Aliases{
		"!foo": "!echo Hello, World!",
	}

	assert.Equal(t, expected, result)
	assert.True(t, called)
}

func TestAliasAddErrors(t *testing.T) {
	cases := map[string]string{
		"missing NAME":  "!alias add",
		"missing VALUE": "!alias add !foo",
	}

	cmd := NewAliasCommand(newMemoryStore(t), ioutil.Discard, func() {})
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestAliasList(t *testing.T) {
	aliases := models.Aliases{
		"!foo": "bar",
		"!bar": "baz",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.AliasesKey, aliases); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewAliasCommand(store, w, func() {})
	if err := runTestApp(cmd, "!alias ls"); err != nil {
		t.Fatal(err)
	}

	for name := range aliases {
		assert.Contains(t, w.String(), name)
	}
}

func TestAliasListErrors(t *testing.T) {
	cmd := NewAliasCommand(newMemoryStore(t), ioutil.Discard, func() {})
	if err := runTestApp(cmd, "!alias ls"); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestAliasRemove(t *testing.T) {
	aliases := models.Aliases{
		"!foo": "bar",
		"!bar": "baz",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.AliasesKey, aliases); err != nil {
		t.Fatal(err)
	}

	var called bool
	invalidate := func() {
		called = true
	}

	cmd := NewAliasCommand(store, ioutil.Discard, invalidate)
	if err := runTestApp(cmd, "!alias rm !foo"); err != nil {
		t.Fatal(err)
	}

	expected := models.Aliases{
		"!bar": "baz",
	}

	result := models.Aliases{}
	if err := store.Read(db.AliasesKey, &result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
	assert.True(t, called)
}

func TestAliasRemoveErrors(t *testing.T) {
	cases := map[string]string{
		"missing NAME":       "!alias rm",
		"NAME doesn't exist": "!alias rm !foo",
	}

	cmd := NewAliasCommand(newMemoryStore(t), ioutil.Discard, func() {})
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestAliasTest(t *testing.T) {
	aliases := models.Aliases{
		"!foo": "{{ .UserID }} in {{ .ChannelID }} says {{ .ArgsString }} (args={{ .Args }})",
	}

	store := newMemoryStore(t)
	if err := store.Write(db.AliasesKey, aliases); err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		Input  string
		Output string
	}{
		"one argument": {
			Input:  "!alias test !foo arg0",
			Output: "user_id in channel_id says arg0 (args=[arg0])",
		},
		"two arguments": {
			Input:  "!alias test !foo arg0 arg1",
			Output: "user_id in channel_id says arg0 arg1 (args=[arg0 arg1])",
		},
		"--user flag": {
			Input:  "!alias test --user UID !foo arg0",
			Output: "UID in channel_id says arg0 (args=[arg0])",
		},
		"--channel flag": {
			Input:  "!alias test --channel CID !foo arg0",
			Output: "user_id in CID says arg0 (args=[arg0])",
		},
	}

	for name := range cases {
		t.Run(name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewAliasCommand(store, w, func() {})
			if err := runTestApp(cmd, cases[name].Input); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, cases[name].Output, w.String())
		})
	}
}

func TestAliasTestErrors(t *testing.T) {
	cases := map[string]string{
		"missing TEXT argument": "!alias test",
		// TODO: Potential addition?
		// "unmatched key":         "!alias keyThatDoesNotExist",
	}

	cmd := NewAliasCommand(newMemoryStore(t), ioutil.Discard, func() {})
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
