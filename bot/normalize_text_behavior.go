package bot

import (
	"strings"

	"github.com/quintilesims/slack"
)

// NewNormalizeTextBehavior returns a behavior that normalizes the text in slack message events
func NewNormalizeTextBehavior() Behavior {
	replacer := strings.NewReplacer("‘", "'", "’", "'", "“", "\"", "”", "\"")
	return func(e slack.RTMEvent) error {
		m, ok := e.Data.(*slack.MessageEvent)
		if !ok {
			return nil
		}

		m.Text = replacer.Replace(m.Text)
		return nil
	}
}
