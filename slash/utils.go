package slash

import (
	"fmt"
	"regexp"
	"strings"
)

// escaped format:  <@U1234|user>
func parseEscapedUser(escaped string) (string, string, error) {
	r := regexp.MustCompile("\\<\\@[a-zA-Z0-9]+|[a-zA-Z0-9]+\\>")
	if !r.MatchString(escaped) {
		return "", "", fmt.Errorf("escaped user is not in valid format")
	}

	// strip '<@' from the front and '>' from the end
	split := strings.SplitN(escaped[2:len(escaped)-1], "|", 2)
	return split[0], split[1], nil
}
