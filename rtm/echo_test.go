package rtm

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	cases := map[string]string{
		"!echo Hello, World!":   "Hello, World!",
		"!echo":                 "",
		"!echo one two three":   "one two three",
		"!echo onetwo    three": "onetwo    three",
	}

	a := NewEchoBehavior()
	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			e := newMessageEvent(input)
			w := bytes.NewBuffer(nil)
			if err := a.OnMessageEvent(e, w); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, w.String())
		})
	}
}

func TestEchoPassthrough(t *testing.T) {
	cases := []string{
		"echo Hello, World!",
		"!ehco Hello, World!",
		"some other command",
		"",
	}

	a := NewEchoBehavior()
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			e := newMessageEvent(c)
			w := bytes.NewBuffer(nil)
			if err := a.OnMessageEvent(e, w); err != nil {
				t.Fatal(err)
			}

			assert.Empty(t, w.String())
		})
	}
}
