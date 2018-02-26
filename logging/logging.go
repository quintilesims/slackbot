package logging

import (
	"io"
	"os"
	"strings"

	"github.com/quintilesims/slackbot/utils"
)

func NewLogWriter(debug bool) io.Writer {
	return utils.WriterFunc(func(p []byte) (n int, err error) {
		if !debug && strings.Contains(string(p), "[DEBUG]") {
			return 0, nil
		}

		return os.Stdout.Write(p)
	})
}
