package lock

type MemoryLock struct {
	isLocked bool
}

func NewMemoryLock() *MemoryLock {
	return &MemoryLock{}
}

func (m *MemoryLock) Lock(wait bool) error {
	if m.isLocked && !wait {
		return NewLockContentionError("Lock is under contention")
	}

	for m.isLocked {
	}
	m.isLocked = true
	return nil
}

func (m *MemoryLock) Unlock() error {
	m.isLocked = false
	return nil
}
