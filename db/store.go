package db

// Store objects are used to read and write data
type Store interface {
	// Keys lists all of the keys in the store
	Keys() ([]string, error)

	// Read will read the value at the specified key into v
	Read(key string, v interface{}) error

	// Write will write v at the specified key
	Write(key string, v interface{}) error
}
