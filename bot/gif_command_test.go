package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGIF(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/search", r.URL.Path)

		query := r.URL.Query()
		fmt.Println(query)
		assert.Equal(t, "key", query.Get("key"))
		assert.Equal(t, "dogs playing poker", query.Get("q"))

		response := TenorSearchResponse{
			Gifs: []Gif{
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
	cmd := NewGIFCommand(server.URL, "key", w)
	if err := runTestApp(cmd, "!gif --explicit dogs playing poker"); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "some url", w.String())
}

func TestGIFErrors(t *testing.T) {
	inputs := []string{
		"!gif",
		"!gif --explicit",
		"!gif --explicit 2 dogs",
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
