package bot

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	cases := map[string]string{
		"!echo":           "",
		"!echo arg0":      "arg0",
		"!echo arg0 arg1": "arg0 arg1",
	}

	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewEchoCommand(w)
			if err := runTestApp(cmd, input); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, w.String())
		})
	}
}
