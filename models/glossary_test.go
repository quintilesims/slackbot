package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlossarySortKeys(t *testing.T) {
	glossary := Glossary{
		"charlie": "",
		"alpha":   "",
		"echo":    "",
		"beta":    "",
		"delta":   "",
	}

	assert.Equal(t, []string{"alpha", "beta", "charlie", "delta", "echo"}, glossary.SortKeys(true))
	assert.Equal(t, []string{"echo", "delta", "charlie", "beta", "alpha"}, glossary.SortKeys(false))
}
