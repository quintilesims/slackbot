package bot

import (
	"github.com/quintilesims/slack"
	"github.com/quintilesims/slackbot/utils"
)

// Behavior executes some functionality when a RTMEvent occurs.
type Behavior func(e slack.RTMEvent) error

// Behaviors adds functionality to a slice of Behavior objects.
type Behaviors []Behavior

// Run will execute each behavior against the event.
// Errors returned from each behavior are aggregated into a single error.
func (behaviors Behaviors) Run(e slack.RTMEvent) error {
	errs := []error{}
	for _, b := range behaviors {
		if err := b(e); err != nil {
			errs = append(errs, err)
		}
	}

	return utils.MultiError(errs)
}
