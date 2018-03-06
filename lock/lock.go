package lock

// Lock objects are used as a point of synchronization 
type Lock interface {
	// Lock will attempt to acquire the Lock
	// If wait is true, the function will block until the lock is next available
	// If wait is false, the function will either acquire the lock or throw a LockContentionError
	Lock(wait bool) error

	// Unlock will release the lock
	Unlock() error
}
