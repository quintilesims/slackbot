package utils

import (
	"crypto/md5"
	"fmt"
	"time"
)

type IDGenerator func() string

func NewGUIDGenerator() IDGenerator {
	return func() string {
		salt := time.Now().Format(time.StampNano)
		return fmt.Sprintf("%x", md5.Sum([]byte(salt)))[:10]
	}
}

func NewStaticIDGenerator(id string) IDGenerator {
	return func() string {
		return id
	}
}
