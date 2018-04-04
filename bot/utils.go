package bot

import (
	"fmt"
	"io"
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

func write(w io.Writer, text string) error {
	if _, err := w.Write([]byte(text)); err != nil {
		return err
	}

	return nil
}

func writef(w io.Writer, format string, tokens ...interface{}) error {
	return write(w, fmt.Sprintf(format, tokens...))
}
