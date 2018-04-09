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

func TestGIFExplicit(t *testing.T) {
	cases := map[string]bool{
		"explicit flag disabled": false,
		"explicit flag enabled":  true,
	}

	for name, explicit := range cases {
		t.Run(name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "/v1/search", r.URL.Path)

				query := r.URL.Query()
				assert.Equal(t, "key", query.Get("key"))
				assert.Equal(t, "dogs playing poker", query.Get("q"))

				expectedSafesearch := "strict"
				if explicit {
					expectedSafesearch = "off"
				}

				assert.Equal(t, expectedSafesearch, query.Get("safesearch"))

				response := TenorSearchResponse{
					Gifs: []Gif{
						{URL: "url"},
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

			input := "!gif dogs playing poker"
			if explicit {
				input = "!gif --explicit dogs playing poker"
			}

			if err := runTestApp(cmd, input); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "url", w.String())
		})
	}
}

func TestGIFLimit(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/search", r.URL.Path)

		query := r.URL.Query()
		assert.Equal(t, "10", query.Get("limit"))

		response := TenorSearchResponse{
			Gifs: []Gif{
				{URL: "url"},
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

	input := "!gif --limit=10 dogs playing poker"
	if err := runTestApp(cmd, input); err != nil {
		t.Fatal(err)
	}
}

func TestGIFErrors(t *testing.T) {
	cases := map[string]string{
		"missing QUERY":           "!gif",
		"missing value for LIMIT": "!gif --limit dog",
	}

	cmd := NewGIFCommand("", "", ioutil.Discard)
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}
