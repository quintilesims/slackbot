package models

import (
	"fmt"
	"sort"
	"strings"
)

// Trivia models hold information about a specific trivia question
type TriviaQuestion struct {
	Question         string
	CorrectAnswer    string
	IncorrectAnswers []string
}

// String will return a string representation of the question and all possible answers
func (t TriviaQuestion) String() string {
	// sort answers in reverse alphabetical order so we display 'True or False?'
	answers := append(t.IncorrectAnswers, t.CorrectAnswer)
	sort.Sort(sort.Reverse(sort.StringSlice(answers)))

	text := fmt.Sprintf("%s\n", t.Question)
	for i, answer := range answers {
		if i == len(answers)-1 {
			text += fmt.Sprintf("or *%s*?", answer)
			break
		}

		text += fmt.Sprintf("*%s*, ", answer)
	}

	// todo: use decoder
	r := strings.NewReplacer("&#039;", "'", "&quot;", "\"")
	return r.Replace(text)
}
