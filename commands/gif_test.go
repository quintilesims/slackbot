package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGif(t *testing.T) {
	w := bytes.NewBuffer(nil)

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/search")

		w.Write([]byte(googleGifHTML))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	client := NewGoogleClient()

	defer server.Close()
	cmd := NewGifCommand(w, client, server.URL)
	if err := runTestApp(cmd, "!gif monkey"); err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("https://media.giphy.com/media/Bl6VoPv34mX2E/200w.gif")
	assert.Equal(t, expected, w.String())
}

func TestGifError(t *testing.T) {
	client := NewGoogleClient()
	w := bytes.NewBuffer(nil)

	cmd := NewGifCommand(w, client, "")
	if err := runTestApp(cmd, "!gif"); err == nil {
		t.Fatal("Error was nil!")
	}
}
