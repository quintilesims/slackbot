package bot

import (
	"fmt"
	"regexp"
)

// date and time layouts
const (
	DateLayout       = "01/02/2006"
	TimeLayout       = "03:04PM"
	DateTimeLayout   = DateLayout + " " + TimeLayout
	DateAtTimeLayout = DateLayout + " at " + TimeLayout
)

func parseEscapedUserID(escaped string) (string, error) {
	// escaped user format: '<@ABC123>'
	r := regexp.MustCompile("\\<\\@.+\\>")
	if !r.MatchString(escaped) {
		return "", fmt.Errorf("Escaped slack user '%s' is not in valid @<username> format", escaped)
	}

	return escaped[2 : len(escaped)-1], nil
}
