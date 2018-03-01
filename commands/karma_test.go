package commands

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/quintilesims/slackbot/common"
	"github.com/quintilesims/slackbot/db"
	"github.com/stretchr/testify/assert"
)

func TestKarma(t *testing.T) {
	karma := map[string]int{
		"dogs":     9999,
		"cats":     -9999,
		"tacos":    50,
		"burritos": 51,
		"pho":      -52,
	}

	store := db.NewMemoryStore()
	if err := store.Write(common.StoreKeyKarma, karma); err != nil {
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
	store := db.NewMemoryStore()
	w := bytes.NewBuffer(nil)
	cmd := NewKarmaCommand(store, w)
	if err := runTestApp(cmd, "!karma"); err == nil {
		t.Fatalf("Error was nil!")
	}
}
