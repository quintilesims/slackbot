package bot

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	cases := map[string]struct {
		Input  string
		Output string
	}{
		"empty string": {
			Input:  "!echo",
			Output: "",
		},
		"one argument": {
			Input:  "!echo arg0",
			Output: "arg0",
		},
		"two argument": {
			Input:  "!echo arg0 arg1",
			Output: "arg0 arg1",
		},
		"flags and commands": {
			Input:  "!echo -v arg0 --stuff arg1",
			Output: "-v arg0 --stuff arg1",
		},
	}

	for name := range cases {
		t.Run(name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewEchoCommand(w)
			if err := runTestApp(cmd, cases[name].Input); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, cases[name].Output, w.String())
		})
	}
}
