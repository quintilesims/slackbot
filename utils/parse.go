package utils

import (
	"fmt"
	"strings"
)

// ParseShell takes a string input and parses it into a []string of arguments.
// Anything wrapped in quotation marks will be treated as a single argument.
// If there is an odd number of quotation marks, an error will be returned.
func ParseShell(input string) ([]string, error) {
	// normalize quotation marks
	r := strings.NewReplacer("‘", "'", "’", "'", "“", "\"", "”", "\"")
	input = r.Replace(input)

	if strings.Count(input, "\"")%2 == 1 {
		return nil, fmt.Errorf("Invalid command: command contains an unpaired quotation mark")
	}

	result := []string{}
	var quoting bool
	var current string
	for _, c := range input {
		switch c {
		case '"':
			quoting = !quoting
			if current != "" {
				result = append(result, current)
			}

			current = ""
		case ' ':
			if quoting {
				current += " "
				continue
			}

			if current != "" {
				result = append(result, current)
			}

			current = ""
		default:
			current += string(c)
		}
	}

	// append the last argument if it exists
	if current != "" {
		result = append(result, current)
	}

	return result, nil
}
