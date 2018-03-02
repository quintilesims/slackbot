package runners

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunnerRun(t *testing.T) {
	var called bool
	r := NewRunner("", func() error {
		called = true
		return nil
	})

	if err := r.Run(); err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
}

func TestRunnerRunEvery(t *testing.T) {
	c := make(chan bool)
	r := NewRunner("", func() error {
		c <- true
		return nil
	})

	ticker := r.RunEvery(time.Nanosecond)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
		select {
		case <-c:
		case <-time.After(time.Millisecond):
			t.Fatal("timeout")
		}
	}
}
