package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCandidateGet(t *testing.T) {
	candidates := Candidates{
		{Name: "alpha"},
		{Name: "beta"},
		{Name: "charlie"},
	}

	cases := map[string]bool{
		"alpha":    true,
		"ALPHA":    true,
		"beta":     true,
		"BETA":     true,
		"alhpa":    false,
		"charlie ": false,
	}

	for name, expected := range cases {
		if _, ok := candidates.Get(name); ok != expected {
			t.Errorf("%s: got %v, expected %v", name, ok, expected)
		}
	}
}

func TestCandidateSort(t *testing.T) {
	candidates := Candidates{
		{Name: "charlie"},
		{Name: "alpha"},
		{Name: "echo"},
		{Name: "beta"},
		{Name: "delta"},
	}

	expected := Candidates{
		{Name: "alpha"},
		{Name: "beta"},
		{Name: "charlie"},
		{Name: "delta"},
		{Name: "echo"},
	}

	candidates.Sort(true)
	assert.Equal(t, expected, candidates)

	expected = Candidates{
		{Name: "echo"},
		{Name: "delta"},
		{Name: "charlie"},
		{Name: "beta"},
		{Name: "alpha"},
	}

	candidates.Sort(false)
	assert.Equal(t, expected, candidates)
}
