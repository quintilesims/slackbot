package bot

import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTextBehavior(t *testing.T) {
	cases := map[string]struct {
		input          string
		expectedOutput string
	}{
		"escape double quotes": {
			"“foo”",
			"\"foo\"",
		},
		"single quotes": {
			"‘foo’",
			"'foo'",
		},
		"sanitize character": {
			"&lt;foo&gt;",
			"<foo>",
		},
	}

	b := NewNormalizeTextBehavior()
	for name := range cases {
		t.Run(name, func(t *testing.T) {
			e := newSlackMessageEvent(cases[name].input)
			if err := b(e); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, cases[name].expectedOutput, e.Data.(*slack.MessageEvent).Text)
		})
	}
}
