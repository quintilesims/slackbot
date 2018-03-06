package lock

import (
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
)

// StoreLock handles locking using a db.Store for storage
type StoreLock struct {
	key   string
	store db.Store
}

// NewStoreLock creates a new StoreLock object with the specified key and store
// The key is used to manage the lock, so multiple instances of StoreLock objects
// can share the same store
func NewStoreLock(key string, store db.Store) *StoreLock {
	return &StoreLock{
		key:   key,
		store: store,
	}
}

func (s *StoreLock) read() (models.Locks, error) {
	locks := models.Locks{}
	if err := s.store.Read(models.StoreKeyLocks, &locks); err != nil {
		return nil, err
	}

	return locks, nil
}

// Lock will attempt to acquire the lock
// If wait is true, the function will block until the lock is next available
// If wait is false, the function will either acquire the lock or throw a LockContentionError
func (s *StoreLock) Lock(wait bool) error {
	// todo: use time multiplier
	for ; ; time.Sleep(time.Second) {
		locks, err := s.read()
		if err != nil {
			return err
		}

		if locks[s.key] && !wait {
			return NewLockContentionError("Lock is under contention")
		}

		if !locks[s.key] {
			locks[s.key] = true
			return s.store.Write(models.StoreKeyLocks, locks)
		}
	}
}

// Unlock will release the lock
func (s *StoreLock) Unlock() error {
	locks, err := s.read()
	if err != nil {
		return err
	}

	locks[s.key] = false
	return s.store.Write(models.StoreKeyLocks, locks)
}
