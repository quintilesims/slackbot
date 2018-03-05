package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func ParseShell(input string) ([]string, error) {
	// normalize quotation marks
	r := strings.NewReplacer("‘", "'", "’", "'", "“", "\"", "”", "\"")
	input = r.Replace(input)

	count := 0
	lastQuote := rune(0)
	f := func(r rune) bool {
		switch {
		case r == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.Is(unicode.Quotation_Mark, r):
			if r == '"' {
				count++
				lastQuote = r
			}
			return false
		default:
			return unicode.IsSpace(r)

		}
	}

	args := strings.FieldsFunc(input, f)
	if count%2 != 0 {
		return nil, fmt.Errorf("Invalid command: command contains an unpaired quotation mark")
	}

	for i := 0; i < len(args); i++ {
		trim := func(r rune) bool { return unicode.Is(unicode.Quotation_Mark, r) }
		args[i] = strings.TrimLeftFunc(args[i], trim)
		args[i] = strings.TrimRightFunc(args[i], trim)
	}

	return args, nil
}

func ParseSlackUser(escaped string) (string, error) {
	// escaped user format: '<@ABC123>'
	r := regexp.MustCompile("\\<\\@[a-zA-Z0-9]+\\>")
	if !r.MatchString(escaped) {
		return "", fmt.Errorf("Invalid user: please enter a valid user by typing `@<username>`")
	}

	return escaped[2 : len(escaped)-1], nil
}
