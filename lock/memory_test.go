package lock

import "testing"

func TestMemoryLock(t *testing.T) {
	testLock(t, NewMemoryLock())
}
