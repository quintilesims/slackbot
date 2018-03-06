package utils

import (
	"io"
	"os"
	"strings"
)

// NewLogWriter will return an io.Writer that writes to stdout
// if debug is false, messages that begin with [DEBUG] will not be written
func NewLogWriter(debug bool) io.Writer {
	return WriterFunc(func(p []byte) (n int, err error) {
		if !debug && strings.Contains(string(p), "[DEBUG]") {
			return 0, nil
		}

		return os.Stdout.Write(p)
	})
}
