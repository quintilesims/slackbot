package bot

import (
	"testing"

	"github.com/quintilesims/slack"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTextBehavior(t *testing.T) {
	cases := map[string]string{
		"foo":         "foo",
		"“foo”":       "\"foo\"",
		"‘foo’":       "'foo'",
		"&lt;foo&gt;": "<foo>",
	}

	b := NewNormalizeTextBehavior()
	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			e := newSlackMessageEvent(input)
			if err := b(e); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, e.Data.(*slack.MessageEvent).Text)
		})
	}
}
