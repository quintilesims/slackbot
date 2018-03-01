package utils

import (
	"crypto/md5"
	"fmt"
	"time"
)

var NewGUID = func(length int) string {
	salt := time.Now().Format(time.StampNano)
	return fmt.Sprintf("%x", md5.Sum([]byte(salt)))[:length]
}
