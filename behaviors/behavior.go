package behaviors

import "github.com/nlopes/slack"

// Behavior executes some functionality when a RTMEvent occurs.
type Behavior func(e slack.RTMEvent) error
