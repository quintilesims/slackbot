package lock

import "errors"

type LockContentionError error

func NewLockContentionError() LockContentionError {
	return LockContentionError(errors.New("lock is under contention"))
}
