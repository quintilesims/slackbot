package behaviors

import (
	"strings"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/common"
	"github.com/quintilesims/slackbot/db"
)

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

		karma := map[string]int{}
		if err := store.Read(common.StoreKeyKarma, &karma); err != nil {
			return err
		}

		// strip '++' or '--'
		key := d.Msg.Text[:len(d.Msg.Text)-2]
		karma[key] = update(karma[key])
		return store.Write(common.StoreKeyKarma, karma)
	}
}
