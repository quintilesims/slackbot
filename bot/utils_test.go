package bot

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	w := bytes.NewBuffer(nil)
	if err := write(w, "Write test"); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Write test", w.String())
}

func TestWritef(t *testing.T) {
	w := bytes.NewBuffer(nil)
	if err := writef(w, "Write test %d", 1); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Write test 1", w.String())
}
