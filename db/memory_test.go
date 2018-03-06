package db

import "testing"

func TestMemoryStore(t *testing.T) {
	testStore(t, NewMemoryStore())
}
