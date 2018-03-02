package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func M(args []string, err error) []interface{} {
	return []interface{}{args, err}
}

type ParseShellOutput struct {
	Args []string
	Err  error
}

func TestParseShell(t *testing.T) {
	cases := map[string]ParseShellOutput{
		"":                         ParseShellOutput{[]string{}, nil},
		"one":                      ParseShellOutput{[]string{"one"}, nil},
		"one two":                  ParseShellOutput{[]string{"one", "two"}, nil},
		"\"one\"":                  ParseShellOutput{[]string{"one"}, nil},
		"'one'":                    ParseShellOutput{[]string{"one"}, nil},
		"one two three":            ParseShellOutput{[]string{"one", "two", "three"}, nil},
		"one \"two\" three":        ParseShellOutput{[]string{"one", "two", "three"}, nil},
		"one \"two three\"":        ParseShellOutput{[]string{"one", "two three"}, nil},
		"one 'two three'":          ParseShellOutput{[]string{"one", "two three"}, nil},
		"\"one two three\"":        ParseShellOutput{[]string{"one two three"}, nil},
		"'one two three'":          ParseShellOutput{[]string{"one two three"}, nil},
		"‘one two three’":          ParseShellOutput{[]string{"one two three"}, nil},
		"“one two three”":          ParseShellOutput{[]string{"one two three"}, nil},
		"one 'two three":           ParseShellOutput{[]string{"one 'two", "three"}, nil},
		"one \"let's have lunch\"": ParseShellOutput{[]string{"one", "let's have lunch"}, nil},
	}

	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, M(expected.Args, expected.Err), M(ParseShell(input)))
		})
	}
}

func TestParseShell_UserInputError(t *testing.T) {
	cases := map[string]ParseShellOutput{
		"one \"two three": ParseShellOutput{[]string{}, fmt.Errorf("Invalid command: command contains an unpaired quotation mark: 'one \"two three'")},
	}

	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			assert.Error(t, expected.Err, M(expected.Args, expected.Err), M(ParseShell(input)))
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
				t.Fatalf("Error was nil")
			}
		})
	}
}
