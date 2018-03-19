package db

import (
	"testing"

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
		AliasesKey,
		InterviewsKey,
		KarmasKey,
	}

	assert.ElementsMatch(t, expected, keys)
}
