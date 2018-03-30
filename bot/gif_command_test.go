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
	cases := []struct {
		name  string
		input string
		url   string
	}{
		{"Clean GIF", "!gif dogs playing poker", "url"},
		{"Explicit GIF", "!gif --explicit dogs playing poker", "saucy url"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "/v1/search", r.URL.Path)

				query := r.URL.Query()
				assert.Equal(t, "key", query.Get("key"))
				assert.Equal(t, "dogs playing poker", query.Get("q"))
				// How do we validate the input?
				// if !c.Bool("explicit") {
				// assert.Equal(t, "strict", query.Get("safesearch"))
				// }

				response := TenorSearchResponse{
					Gifs: []Gif{
						{URL: c.url},
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
			if err := runTestApp(cmd, c.input); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.url, w.String())
		})
	}
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
