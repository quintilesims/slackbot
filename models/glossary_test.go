package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlossarySortKeys(t *testing.T) {
	glossary := Glossary{
		"foo": "bar",
		"bar": "baz",
		"baz": "foo",
	}

	assert.Equal(t, []string{"bar", "baz", "foo"}, glossary.SortKeys(true))
	assert.Equal(t, []string{"foo", "baz", "bar"}, glossary.SortKeys(false))
}
