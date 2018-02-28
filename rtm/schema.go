package rtm

import (
	"io"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/utils"
)

// todo: should this have a Validate() error function?
type BehaviorSchema struct {
	Name           string
	Usage          string
	Help           string
	Init           func() error
	OnMessageEvent func(e *slack.MessageEvent, w io.Writer) error
}

type Behaviors []*BehaviorSchema

func (bh Behaviors) Do(f func(*BehaviorSchema) error) error {
	errs := []error{}
	for _, b := range bh {
		if err := f(b); err != nil {
			errs = append(errs, err)
		}
	}

	return utils.MultiError(errs)
}

func (bh Behaviors) Init() error {
	return bh.Do(func(b *BehaviorSchema) error {
		if b.Init != nil {
			return b.Init()
		}

		return nil
	})
}

func (bh Behaviors) OnMessageEvent(e *slack.MessageEvent, w io.Writer) error {
	return bh.Do(func(b *BehaviorSchema) error {
		if b.OnMessageEvent != nil {
			return b.OnMessageEvent(e, w)
		}

		return nil
	})
}
