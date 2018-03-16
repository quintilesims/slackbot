package bot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGIF(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/gifs/search", r.URL.Path)

		query := r.URL.Query()
		assert.Equal(t, "token", query.Get("api_key"))
		assert.Equal(t, "pg", query.Get("rating"))
		assert.Equal(t, "dogs playing poker", query.Get("q"))

		response := GiphySearchResponse{
			Gifs: []GiphyGif{
				{URL: "some url"},
			},
		}

		b, err := json.Marshal(response)
		if err != nil {
			t.Fatal(err)
		}

		w.Write(b)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	w := bytes.NewBuffer(nil)
	cmd := NewGIFCommand(server.URL, "token", w)
	if err := runTestApp(cmd, "!gif --rating pg dogs playing poker"); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "some url", w.String())
}

func TestGIFErrors(t *testing.T) {
	inputs := []string{
		"!gif",
		"!gif --rating r",
		"!gif --rating 2 dogs",
	}

	cmd := NewGIFCommand("", "", ioutil.Discard)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
