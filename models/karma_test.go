package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKarmaSortKeys(t *testing.T) {
	karmas := Karmas{
		"three": {Upvotes: 5, Downvotes: 5},
		"one":   {Upvotes: 0, Downvotes: 10},
		"five":  {Upvotes: 10, Downvotes: 0},
		"two":   {Upvotes: 0, Downvotes: 5},
		"four":  {Upvotes: 5, Downvotes: 0},
	}

	assert.Equal(t, []string{"five", "four", "three", "two", "one"}, karmas.SortKeys(true))
	assert.Equal(t, []string{"one", "two", "three", "four", "five"}, karmas.SortKeys(false))
}
