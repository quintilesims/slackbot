package rtm

import (
	"bytes"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	w := bytes.NewBuffer(nil)
	e := newMessageEvent("!echo Hello, World!")
	if err := NewEchoAction().OnMessageEvent(e, w); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Hello, World!", w.String())
}

func TestEchoPassthrough(t *testing.T) {
	w := bytes.NewBuffer(nil)
	e := &slack.MessageEvent{Msg: slack.Msg{Text: "!echo Hello, World!"}}
	if err := NewEchoAction().OnMessageEvent(e, w); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Hello, World!", w.String())
}
