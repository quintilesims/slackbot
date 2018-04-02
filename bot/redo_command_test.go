package bot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedo(t *testing.T) {
	var called bool
	trigger := func() error {
		called = true
		return nil
	}

	cmd := NewRedoCommand(trigger)
	if err := runTestApp(cmd, "!redo"); err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
}

func TestRedoError(t *testing.T) {
	trigger := func() error {
		return fmt.Errorf("some error")
	}

	cmd := NewRedoCommand(trigger)
	if err := runTestApp(cmd, "!redo"); err == nil {
		t.Fatal("Error was nil!")
	}
}
