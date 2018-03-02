package lock

import "fmt"

type LockContentionError error

func NewLockContentionError(format string, tokens ...interface{}) LockContentionError {
	return LockContentionError(fmt.Errorf(format, tokens...))
}
