package utils

import (
	"crypto/md5"
	"fmt"
	"time"
)

// IDGenerator is a function that generates a string
type IDGenerator func() string

// NewGUIDGenerator returns an IDGenerator that generates a globally unique identifier
func NewGUIDGenerator() IDGenerator {
	return func() string {
		salt := time.Now().Format(time.StampNano)
		return fmt.Sprintf("%x", md5.Sum([]byte(salt)))[:10]
	}
}

// NewStaticIDGenerator returns an IDGenerator that only generates the specified id
func NewStaticIDGenerator(id string) IDGenerator {
	return func() string {
		return id
	}
}
