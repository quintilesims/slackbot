package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseShell(t *testing.T) {
	t.Skip("TODO: Fix")

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

func TestParseSlackUser(t *testing.T) {
	cases := map[string]string{
		"<@ABC123>": "ABC123",
		"<@j3lfIa>": "j3lfIa",
	}

	for escaped, expected := range cases {
		t.Run(escaped, func(t *testing.T) {
			result, err := ParseSlackUser(escaped)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, result)
		})
	}
}

func TestParseSlackUserError(t *testing.T) {
	inputs := []string{
		"@ABC123",
		"username",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if _, err := ParseSlackUser(input); err == nil {
				t.Fatal("Error was nil")
			}
		})
	}
}
