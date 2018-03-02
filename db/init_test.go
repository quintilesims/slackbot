package db

import (
	"testing"

	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestInitKarmaStore(t *testing.T) {
	store := NewMemoryStore()
	if err := initKarmaStore(store); err != nil {
		t.Fatal(err)
	}

	keys, err := store.Keys()
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, keys, models.StoreKeyKarma)
}

func TestInitRemindersStore(t *testing.T) {
	store := NewMemoryStore()
	if err := initRemindersStore(store); err != nil {
		t.Fatal(err)
	}

	keys, err := store.Keys()
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, keys, models.StoreKeyReminders)
}
