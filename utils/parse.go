package utils

import (
	"strings"
	"unicode"
)

func ParseShell(input string) []string {
	// normalize quotation marks
	r := strings.NewReplacer("‘", "'", "’", "'", "“", "\"", "”", "\"")
	input = r.Replace(input)

	lastQuote := rune(0)
	f := func(r rune) bool {
		switch {
		case r == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.Is(unicode.Quotation_Mark, r):
			lastQuote = r
			return false
		default:
			return unicode.IsSpace(r)

		}
	}

	args := strings.FieldsFunc(input, f)
	for i := 0; i < len(args); i++ {
		trim := func(r rune) bool { return unicode.Is(unicode.Quotation_Mark, r) }
		args[i] = strings.TrimLeftFunc(args[i], trim)
		args[i] = strings.TrimRightFunc(args[i], trim)
	}

	return args
}
