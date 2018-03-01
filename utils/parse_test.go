package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseShell(t *testing.T) {
	cases := map[string][]string{
		"":                  []string{},
		"one":               []string{"one"},
		"one two":           []string{"one", "two"},
		"\"one\"":           []string{"one"},
		"'one'":             []string{"one"},
		"one two three":     []string{"one", "two", "three"},
		"one \"two\" three": []string{"one", "two", "three"},
		"one \"two three\"": []string{"one", "two three"},
		"one 'two three'":   []string{"one", "two three"},
		"\"one two three\"": []string{"one two three"},
		"'one two three'":   []string{"one two three"},
		"‘one two three’":   []string{"one two three"},
		"“one two three”":   []string{"one two three"},
		"one 'two three":    []string{"one", "'two", "three"},
		"one \"two three":   []string{"one", "\"two", "three"},
	}

	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, expected, ParseShell(input))
		})
	}
}
