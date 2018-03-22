package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseShell(t *testing.T) {
	cases := map[string][]string{
		"":                           []string{},
		"one":                        []string{"one"},
		"one two":                    []string{"one", "two"},
		"\"one\"":                    []string{"one"},
		"'one'":                      []string{"'one'"},
		"one two three":              []string{"one", "two", "three"},
		"one \"two\" three":          []string{"one", "two", "three"},
		"one \"two three\"":          []string{"one", "two three"},
		"one 'two three'":            []string{"one", "'two", "three'"},
		"\"one two three\"":          []string{"one two three"},
		"'one two three'":            []string{"'one", "two", "three'"},
		"one 'two three":             []string{"one", "'two", "three"},
		"one\" two\" three":          []string{"one", " two", "three"},
		"one \"two \"three\" four\"": []string{"one", "two ", "three", " four"},
	}

	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			out, err := ParseShell(input)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, out)
		})
	}
}

func TestParseShellErrors(t *testing.T) {
	inputs := []string{
		"one \"two three",
		"one \"\"\"two three",
		"one\" two\"\" three",
		"one \"two \"three\"\" four\"",
		"one '\"two' three",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if _, err := ParseShell(input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
