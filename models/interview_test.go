package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInterviewEquals(t *testing.T) {
	cases := []struct {
		Name     string
		First    Interview
		Second   Interview
		Expected bool
	}{
		{
			Name:     "empty",
			Expected: true,
		},
		{
			Name: "same name, different time",
			First: Interview{
				Candidate: "name",
				Time:      time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			Second: Interview{
				Candidate: "name",
				Time:      time.Date(0, 0, 0, 1, 0, 0, 0, time.UTC),
			},
			Expected: false,
		},
		{
			Name: "same time, different name",
			First: Interview{
				Candidate: "name1",
				Time:      time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			Second: Interview{
				Candidate: "name2",
				Time:      time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			Expected: false,
		},
		{
			Name: "same name and date",
			First: Interview{
				Candidate: "name",
				Time:      time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			Second: Interview{
				Candidate: "name",
				Time:      time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			Expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert.Equal(t, c.Expected, c.First.Equals(c.Second))
		})
	}
}

func TestInterviewsSort(t *testing.T) {
	now := time.Now()
	interviews := Interviews{
		{Candidate: "Today", Time: now},
		{Candidate: "Yesterday", Time: now.AddDate(0, 0, -1)},
		{Candidate: "Tomorrow", Time: now.AddDate(0, 0, 1)},
	}

	expected := Interviews{
		{Candidate: "Yesterday", Time: now.AddDate(0, 0, -1)},
		{Candidate: "Today", Time: now},
		{Candidate: "Tomorrow", Time: now.AddDate(0, 0, 1)},
	}

	interviews.Sort(true)
	assert.Equal(t, expected, interviews)

	expected = Interviews{
		{Candidate: "Tomorrow", Time: now.AddDate(0, 0, 1)},
		{Candidate: "Today", Time: now},
		{Candidate: "Yesterday", Time: now.AddDate(0, 0, -1)},
	}

	interviews.Sort(false)
	assert.Equal(t, expected, interviews)

}
