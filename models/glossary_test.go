package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlossarySortKeys(t *testing.T) {
	candidates := Candidates{
		"charlie": nil,
		"alpha":   nil,
		"echo":    nil,
		"beta":    nil,
		"delta":   nil,
	}

	assert.Equal(t, []string{"alpha", "beta", "charlie", "delta", "echo"}, candidates.SortKeys(true))
	assert.Equal(t, []string{"echo", "delta", "charlie", "beta", "alpha"}, candidates.SortKeys(false))
}
