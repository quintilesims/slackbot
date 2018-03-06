package utils

import (
	"fmt"
	"regexp"
	"strings"
)

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
			} else {
				if current != "" {
					result = append(result, current)
				}

				current = ""
			}
		default:
			current += string(c)
		}
	}

	// End of string, append the last argument if it exists
	if current != "" {
		result = append(result, current)
	}

	return result, nil
}

func ParseSlackUser(escaped string) (string, error) {
	// escaped user format: '<@ABC123>'
	r := regexp.MustCompile("\\<\\@[a-zA-Z0-9]+\\>")
	if !r.MatchString(escaped) {
		return "", fmt.Errorf("Invalid user: please enter a valid user by typing `@<username>`")
	}

	return escaped[2 : len(escaped)-1], nil
}
