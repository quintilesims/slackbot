package models

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/utils"
)

// AliasContext objects are passed into an alias's value template when the alias is executed
type AliasContext struct {
	ChannelID  string
	UserID     string
	Args       []string
	ArgsString string
}

// The Aliases object is used to manage aliases in a db.Store
type Aliases map[string]string

// Apply will update the MessageEvent's text if it matches any of the Aliases' patterns
func (a Aliases) Apply(m *slack.MessageEvent) error {
	args, err := utils.ParseShell(m.Text)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return nil
	}

	value, ok := a[args[0]]
	if !ok {
		return nil
	}

	t, err := template.New("").Parse(value)
	if err != nil {
		return err
	}

	context := AliasContext{
		ChannelID:  m.Channel,
		UserID:     m.User,
		Args:       args[1:],
		ArgsString: strings.Join(args[1:], " "),
	}

	b := bytes.NewBuffer(nil)
	if err := t.Execute(b, context); err != nil {
		return err
	}

	m.Text = b.String()
	return nil
}
