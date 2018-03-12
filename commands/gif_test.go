package commands

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	cache "github.com/zpatrick/go-cache"
)

func TestGif(t *testing.T) {
	w := bytes.NewBuffer(nil)
	c := cache.New()

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/search")

		w.Write([]byte(googleGifHTML))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	client := NewGoogleClient()

	defer server.Close()
	cmd := NewGifCommand(w, client, c, server.URL)
	if err := runTestApp(cmd, "!gif monkey"); err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("https://media.giphy.com/media/Bl6VoPv34mX2E/200w.gif")
	assert.Equal(t, expected, w.String())
}

func TestGifError(t *testing.T) {
	client := NewGoogleClient()
	w := bytes.NewBuffer(nil)
	c := cache.New()

	cmd := NewGifCommand(w, client, c, "")
	if err := runTestApp(cmd, "!gif"); err == nil {
		t.Fatal("Error was nil!")
	}
}

func TestGif_Cache(t *testing.T) {
	w := bytes.NewBuffer(nil)
	c := cache.New()

	gif := &Gif{ContentURL: "https://media.giphy.com/media/Bl6VoPv34mX2E/200w.gif"}
	c.Add("monkey", gif)

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/search")

		w.Write([]byte(googleGifHTML))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	client := NewGoogleClient()

	defer server.Close()
	cmd := NewGifCommand(w, client, c, server.URL)
	if err := runTestApp(cmd, "!gif monkey"); err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("https://media.giphy.com/media/Bl6VoPv34mX2E/200w.gif")
	assert.Equal(t, expected, w.String())
}

func TestGif_Cache_Error(t *testing.T) {
	w := bytes.NewBuffer(nil)
	c := cache.New()

	c.Add("monkey", "https://media.giphy.com/media/Bl6VoPv34mX2E/200w.gif")

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/search")

		w.Write([]byte(googleGifHTML))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	client := NewGoogleClient()

	defer server.Close()
	cmd := NewGifCommand(w, client, c, server.URL)
	if err := runTestApp(cmd, "!gif monkey"); err != nil {
		assert.Equal(t, err, errors.New("Cache error"))
	}
}
