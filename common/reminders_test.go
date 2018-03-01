package common

import (
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/stretchr/testify/assert"
)

func TestInitRemindersStore(t *testing.T) {
	store := db.NewMemoryStore()
	if err := initRemindersStore(store); err != nil {
		t.Fatal(err)
	}

	keys, err := store.Keys()
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, keys, StoreKeyReminders)
}
