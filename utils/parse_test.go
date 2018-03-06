package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseShell(t *testing.T) {
	t.Skip("TODO: Fix")

	cases := map[string][]string{
		"":                  {},
		"one":               {"one"},
		"one two":           {"one", "two"},
		"\"one\"":           {"one"},
		"'one'":             {"one"},
		"one two three":     {"one", "two", "three"},
		"one \"two\" three": {"one", "two", "three"},
		"one \"two three\"": {"one", "two three"},
		"one 'two three'":   {"one", "two three"},
		"\"one two three\"": {"one two three"},
		"'one two three'":   {"one two three"},
		"‘one two three’":   {"one two three"},
		"“one two three”":   {"one two three"},
		"one 'two three":    {"one", "'two", "three"},
		"one \"two three":   {"one", "\"two", "three"},
	}

	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, expected, ParseShell(input))
		})
	}
}
