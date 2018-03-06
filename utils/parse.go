package utils

import (
	"strings"
)

// ParseShell takes a string input and parses it into a []string of arguments
// anything wrapped in quotation marks will be treated as a single object
// if there is an odd number of quotation marks, an error will be returned
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
