package lock

import "fmt"

// LockContentionError occurs when Lock.Lock(false) is called and the Lock is currently acquired
type LockContentionError error

// NewLockContentionError creates a new LockContentError object
func NewLockContentionError(format string, tokens ...interface{}) LockContentionError {
	return LockContentionError(fmt.Errorf(format, tokens...))
}
