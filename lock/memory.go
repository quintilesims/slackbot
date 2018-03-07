package lock

// MemoryLock handles locking using in-memory storage
type MemoryLock struct {
	isLocked bool
}

// NewMemoryLock creates a new MemoryLock object
func NewMemoryLock() *MemoryLock {
	return &MemoryLock{}
}

// Lock will attempt to acquire the lock.
// If wait is true, the function will block until the lock is released.
// If wait is false, the function will either acquire the lock or throw a LockContentionError.
func (m *MemoryLock) Lock(wait bool) error {
	if m.isLocked && !wait {
		return NewLockContentionError("Lock is under contention")
	}

	for m.isLocked {
	}

	m.isLocked = true
	return nil
}

// Unlock will release the lock
func (m *MemoryLock) Unlock() error {
	m.isLocked = false
	return nil
}
