package bot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestTriviaAnswer(t *testing.T) {
	question := models.TriviaQuestion{
		Question:         "What is the best movie of all time?",
		CorrectAnswer:    "White Chicks",
		IncorrectAnswers: []string{"Casablanca", "Citizen Kane"},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.TriviaKey, question); err != nil {
		t.Fatal(err)
	}

	cases := map[string]bool{
		"!trivia answer white chicks": true,
		"!trivia answer White Chicks": true,
		"!trivia answer casablanca":   false,
		"!trivia answer Citizen Kane": false,
		"!trivia answer foo":          false,
	}

	for input, isCorrect := range cases {
		t.Run(input, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			cmd := NewTriviaCommand(store, "", w)
			if err := runTestApp(cmd, input); err != nil {
				t.Fatal(err)
			}

			expected := "is not the correct answer"
			if isCorrect {
				expected = "is the correct answer"
			}

			assert.Contains(t, w.String(), expected)
		})
	}
}

func TestTriviaAnswerErrors(t *testing.T) {
	inputs := []string{
		"!trivia answer",
		"!trivia answer foo",
	}

	store := newMemoryStore(t)
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			cmd := NewTriviaCommand(store, "", ioutil.Discard)
			if err := runTestApp(cmd, input); err == nil {
				t.Fatal("Error was nil!")
			}
		})
	}
}

func TestTriviaNew(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/", r.URL.Path)

		query := r.URL.Query()
		assert.Equal(t, "1", query.Get("amount"))
		assert.Equal(t, "hard", query.Get("difficulty"))

		// opentdb send escaped html
		response := TriviaResponse{
			Questions: []TriviaQuestion{
				{
					Question:         "Which is the world&#39;s &quot;greatest&quot; band of all time?",
					CorrectAnswer:    "Hoobastank",
					IncorrectAnswers: []string{"Smashmouth", "Matchbox Twenty", "My Chemical Romance"},
				},
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

	store := newMemoryStore(t)
	w := bytes.NewBuffer(nil)
	cmd := NewTriviaCommand(store, server.URL, w)
	if err := runTestApp(cmd, "!trivia new --difficulty hard"); err != nil {
		t.Fatal(err)
	}

	expected := models.TriviaQuestion{
		Question:         "Which is the world's \"greatest\" band of all time?",
		CorrectAnswer:    "Hoobastank",
		IncorrectAnswers: []string{"Smashmouth", "Matchbox Twenty", "My Chemical Romance"},
	}

	assert.Contains(t, w.String(), expected.String())

	result := models.TriviaQuestion{}
	if err := store.Read(db.TriviaKey, &result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, result)
}

func TestTriviaShow(t *testing.T) {
	question := models.TriviaQuestion{
		Question:         "Complete the sentence: Weezy F. Baby and the 'F' is for ...",
		CorrectAnswer:    "phenomenal",
		IncorrectAnswers: []string{"fear", "flaw", "failing"},
	}

	store := newMemoryStore(t)
	if err := store.Write(db.TriviaKey, question); err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBuffer(nil)
	cmd := NewTriviaCommand(store, "", w)
	if err := runTestApp(cmd, "!trivia show"); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, w.String(), question.String())
}

func TestTriviaShowError(t *testing.T) {
	cmd := NewTriviaCommand(newMemoryStore(t), "", ioutil.Discard)
	if err := runTestApp(cmd, "!trivia show"); err == nil {
		t.Fatal("Error was nil!")
	}
}
