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

func TestDefine(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/words", r.URL.Path)

		query := r.URL.Query()
		assert.Equal(t, "1", query.Get("max"))
		assert.Equal(t, "d", query.Get("md"))
		assert.Equal(t, "ice cream", query.Get("sp"))

		response := DatamuseResponse{
			{
				Definitions: []string{
					"frozen dessert containing cream and sugar and flavoring",
				},
				Word: "ice cream",
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
	cmd := NewDefineCommand(server.URL, w)

	if err := runTestApp(cmd, "!define ice cream"); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, w.String(), "frozen dessert containing cream and sugar and flavoring")
}

func TestDefineErrors(t *testing.T) {
	cmd := NewDefineCommand("", ioutil.Discard)

	if err := runTestApp(cmd, "!define"); err == nil {
		t.Fatal("Error was nil!")
	}
}
