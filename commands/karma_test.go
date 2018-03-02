package commands

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestKarma(t *testing.T) {
	karma := models.Karma{
		"dogs":     9999,
		"cats":     -9999,
		"tacos":    50,
		"burritos": 51,
		"pho":      -52,
	}

	store := newMemoryStore(t)
	if err := store.Write(models.StoreKeyKarma, karma); err != nil {
		t.Fatal(err)
	}

	for k, v := range karma {
		t.Run(k, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewKarmaCommand(store, w)
			if err := runTestApp(cmd, "!karma %s", k); err != nil {
				t.Fatal(err)
			}

			expected := fmt.Sprintf("karma for '%s': %d", k, v)
			assert.Equal(t, expected, w.String())
		})
	}
}

func TestKarmaError(t *testing.T) {
	store := newMemoryStore(t)
	w := bytes.NewBuffer(nil)
	cmd := NewKarmaCommand(store, w)
	if err := runTestApp(cmd, "!karma"); err == nil {
		t.Fatalf("Error was nil!")
	}
}
