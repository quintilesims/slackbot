package bot

import (
	"fmt"
	"io"
	"regexp"

	"github.com/quintilesims/slack"
)

// date and time layouts
const (
	DateLayout       = "01/02"
	TimeLayout       = "03:04PM"
	DateTimeLayout   = "01/0203:04PM"
	DateAtTimeLayout = "01/02 at 03:04PM"
)

func parseSlackUser(client slack.SlackClient, escaped string) (*slack.User, error) {
	// escaped user format: '<@ABC123>'
	r := regexp.MustCompile("\\<\\@.+\\>")
	if !r.MatchString(escaped) {
		return nil, fmt.Errorf("Escaped slack user '%s' is not in valid @<username> format", escaped)
	}

	userID := escaped[2 : len(escaped)-1]
	return client.GetUserInfo(userID)
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
