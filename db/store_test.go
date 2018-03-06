package db

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testStore(t *testing.T, store Store) {
	writes := map[string]interface{}{
		"k1": 1,
		"k2": "two",
		"k3": []string{"one", "two", "three"},
		"k4": map[string]int{"one": 1},
	}

	for k, v := range writes {
		if err := store.Write(k, v); err != nil {
			t.Fatal(err)
		}
	}

	keys, err := store.Keys()
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, keys, len(writes))
	for _, k := range keys {
		var v interface{}
		if err := store.Read(k, &v); err != nil {
			t.Fatal(err)
		}

		// use string comparison as object literal comparisons are unreliable
		assert.Equal(t, fmt.Sprintf("%v", writes[k]), fmt.Sprintf("%v", v))
	}

	key := strconv.Itoa(rand.Int())
	if err, ok := store.Read(key, nil).(MissingEntryError); !ok {
		t.Errorf("Error was not MissingEntryError: %#v", err)
	}
}
