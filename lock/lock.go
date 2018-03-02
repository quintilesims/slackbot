package lock

type Lock interface {
	Lock(wait bool) error
	Unlock() error
}
