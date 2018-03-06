package behaviors

import (
	"strings"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
)

// NewKarmaTrackingBehavior will update karma in the provided store
// Karma updates are triggered by the presence of '++' or '--' at the end of a message
func NewKarmaTrackingBehavior(store db.Store) Behavior {
	return func(e slack.RTMEvent) error {
		d, ok := e.Data.(*slack.MessageEvent)
		if !ok {
			return nil
		}

		var update func(i int) int
		switch {
		case strings.HasSuffix(d.Msg.Text, "++"):
			update = func(i int) int { return i + 1 }
		case strings.HasSuffix(d.Msg.Text, "--"):
			update = func(i int) int { return i - 1 }
		default:
			return nil
		}

		karma := models.Karma{}
		if err := store.Read(models.StoreKeyKarma, &karma); err != nil {
			return err
		}

		// strip '++' or '--'
		key := d.Msg.Text[:len(d.Msg.Text)-2]
		karma[key] = update(karma[key])
		return store.Write(models.StoreKeyKarma, karma)
	}
}
