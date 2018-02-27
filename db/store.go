package db

type Store interface {
	Keys() ([]string, error)
	Read(key string, v interface{}) error
	Write(key string, v interface{}) error
}
