package lock

import (
	"testing"

	"github.com/quintilesims/slackbot/db"
)

func TestStoreLock(t *testing.T) {
	store := db.NewMemoryStore()
	if err := db.Init(store); err != nil {
		t.Fatal(err)
	}

	testLock(t, NewStoreLock("key", store))
}
