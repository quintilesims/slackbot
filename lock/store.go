package lock

import (
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
)

type StoreLock struct {
	key   string
	store db.Store
}

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

func (s *StoreLock) Unlock() error {
	locks, err := s.read()
	if err != nil {
		return err
	}

	locks[s.key] = false
	return s.store.Write(models.StoreKeyLocks, locks)
}
