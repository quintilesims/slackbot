package models

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/quintilesims/slack"
	glob "github.com/ryanuber/go-glob"
)

// Alias models hold information about a specific alias
type Alias struct {
	Pattern  string
	Template string
}

// Apply will update the MessageEvent's text if it matches the Alias's pattern
func (a Alias) Apply(m *slack.MessageEvent) error {
	if !glob.Glob(a.Pattern, m.Text) {
		return nil
	}

	funcMap := template.FuncMap{
		"replace": func(input, from, to string) string {
			return strings.Replace(input, from, to, -1)
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).Parse(a.Template)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	if err := tmpl.Execute(b, m); err != nil {
		return nil
	}

	m.Text = b.String()
	return nil
}

// The Aliases object is used to manage Alias instances in a db.Store
type Aliases map[string]Alias
