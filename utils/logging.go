package utils

import (
	"io"
	"os"
	"strings"
)

func NewLogWriter(debug bool) io.Writer {
	return WriterFunc(func(p []byte) (n int, err error) {
		if !debug && strings.Contains(string(p), "[DEBUG]") {
			return 0, nil
		}

		return os.Stdout.Write(p)
	})
}
