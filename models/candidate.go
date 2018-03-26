package models

import "sort"

// Candidate models hold information about a specific candidate
type Candidate map[string]string

// The Candidates object is used to manage Candidate instances in a db.Store
// by using the Candidates' names as keys
type Candidates map[string]Candidate

// SortKeys will return a slice of ordered keys.
// If ascending is true, keys are returned in alphabetical order.
// If ascending is false, keys are returned in reverse alphabetical order.
func (c Candidates) SortKeys(ascending bool) []string {
	keys := make(sort.StringSlice, 0, len(c))
	for key := range c {
		keys = append(keys, key)
	}

	if ascending {
		sort.Sort(keys)
	} else {
		sort.Sort(sort.Reverse(keys))
	}

	return keys
}
