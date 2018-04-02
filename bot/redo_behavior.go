package bot

import (
	"fmt"
	"strings"

	"github.com/quintilesims/slack"
	cache "github.com/zpatrick/go-cache"
)

// RedoBehavior tracks and executes slack.MessageEvent instances.
// Events are recorded using the Record() function,
// and the last event can be executed by using the Trigger() function.
type RedoBehavior struct {
	eventChan  chan slack.RTMEvent
	eventCache *cache.Cache
}

// NewRedoBehavior will create a new instance of a RedoBehavior.
func NewRedoBehavior(c chan slack.RTMEvent) *RedoBehavior {
	return &RedoBehavior{
		eventChan:  c,
		eventCache: cache.New(),
	}
}

// Record will record e as the last event for the specified channel.
// Events cannot be a *slack.MessageEvent whose text starts with "!redo".
func (r *RedoBehavior) Record(channelID string, e slack.RTMEvent) error {
	if m, ok := e.Data.(*slack.MessageEvent); ok && strings.HasPrefix(m.Text, "!redo") {
		return fmt.Errorf("Cannot record MessageEvent starting with !redo")
	}

	r.eventCache.Set(channelID, e)
	return nil
}

// Trigger will send the last event on the specified channel to the RedoBehavior's RTMEvent channel.
// An error will be thrown if no event is recorded for the specified channel.
func (r *RedoBehavior) Trigger(channelID string) error {
	e, ok := r.eventCache.Getf(channelID)
	if !ok {
		return fmt.Errorf("No event recorded for channel %s", channelID)
	}

	go func() { r.eventChan <- e.(slack.RTMEvent) }()
	return nil
}
