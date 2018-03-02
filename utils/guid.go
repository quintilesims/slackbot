package utils

import (
	"crypto/md5"
	"fmt"
	"time"
)

func NewGUID() string {
	salt := time.Now().Format(time.StampNano)
	return fmt.Sprintf("%x", md5.Sum([]byte(salt)))[:10]
}
