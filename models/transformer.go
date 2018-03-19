package models

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/quintilesims/slack"
	glob "github.com/ryanuber/go-glob"
)

// Transformer models hold information about a specific transformation
type Transformer struct {
	Pattern  string
	Template string
}

func (t Transformer) Apply(m *slack.Message) error {
	if !glob.Glob(t.Pattern, m.Text) {
		return nil
	}

	funcMap := template.FuncMap{
		"replace": func(input, from, to string) string {
			return strings.Replace(input, from, to, -1)
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).Parse(t.Template)
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

// The Transformers object is used to manage Transformer instances in a db.Store
type Transformers map[string]Transformer
