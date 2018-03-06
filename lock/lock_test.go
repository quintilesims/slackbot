package lock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testLock(t *testing.T, lock Lock) {
	if err := lock.Lock(false); err != nil {
		t.Fatal(err)
	}

	if err, ok := lock.Lock(false).(LockContentionError); !ok {
		t.Fatalf("Error was not LockContentionError: %#v", err)
	}

	c := make(chan error)
	go func() {
		c <- lock.Lock(true)
	}()

	if err := lock.Unlock(); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-c:
		assert.Nil(t, err)
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}
