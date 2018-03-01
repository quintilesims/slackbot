package common

import (
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/stretchr/testify/assert"
)

func TestInitKarmaStore(t *testing.T) {
	store := db.NewMemoryStore()
	if err := initKarmaStore(store); err != nil {
		t.Fatal(err)
	}

	keys, err := store.Keys()
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, keys, StoreKeyKarma)
}
