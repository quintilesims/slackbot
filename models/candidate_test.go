package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCandidateSortKeys(t *testing.T) {
	karmas := Candidates{
		"charlie": nil,
		"alpha":   nil,
		"echo":    nil,
		"beta":    nil,
		"delta":   nil,
	}

	assert.Equal(t, []string{"alpha", "beta", "charlie", "delta", "echo"}, karmas.SortKeys(true))
	assert.Equal(t, []string{"echo", "delta", "charlie", "beta", "alpha"}, karmas.SortKeys(false))
}
