package commands

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGif(t *testing.T) {
	w := bytes.NewBuffer(nil)
	cmd := NewGifCommand(w)
	if err := runTestApp(cmd, "!gif monkey"); err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("gif url")
	assert.Equal(t, expected, w.String())
}

func TestGifError(t *testing.T) {
	w := bytes.NewBuffer(nil)
	cmd := NewGifCommand(w)
	if err := runTestApp(cmd, "!gif"); err == nil {
		t.Fatal("Error was nil!")
	}
}
