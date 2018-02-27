package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	b := bytes.NewBuffer(nil)
	cmd := NewEchoCommand(b)
	if err := cmd.run([]string{"Hello,", "World!"}); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Hello, World!", b.String())
}

func TestEchoError(t *testing.T) {
	b := bytes.NewBuffer(nil)
	cmd := NewEchoCommand(b)
	if err := cmd.run([]string{}); err == nil {
		t.Fatalf("Error was nil!")
	}
}
