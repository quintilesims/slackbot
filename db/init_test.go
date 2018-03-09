package db

import (
	"testing"

	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	store := NewMemoryStore()
	if err := Init(store); err != nil {
		t.Fatal(err)
	}

	keys, err := store.Keys()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		models.StoreKeyCallbacks,
		models.StoreKeyChecklists,
		models.StoreKeyInterviews,
		models.StoreKeyLocks,
		models.StoreKeyReminders,
	}

	assert.ElementsMatch(t, expected, keys)
}
