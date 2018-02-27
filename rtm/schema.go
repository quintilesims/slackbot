package rtm

import (
	"io"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/utils"
)

type ActionSchema struct {
	Name           string
	Usage          string
	Help           string
	Init           func() error
	OnMessageEvent func(e *slack.MessageEvent, w io.Writer) error
}

type Actions []*ActionSchema

func (ac Actions) Init() error {
	errs := []error{}
	for _, a := range ac {
		if a.Init != nil {
			if err := a.Init(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utils.MultiError(errs)
}

func (ac Actions) OnMessageEvent(e *slack.MessageEvent, w io.Writer) error {
	errs := []error{}
	for _, a := range ac {
		if a.OnMessageEvent != nil {
			if err := a.OnMessageEvent(e, w); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utils.MultiError(errs)
}
