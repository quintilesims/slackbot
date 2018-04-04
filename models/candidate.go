package models

import (
	"sort"
	"strings"
)

// Candidate models hold information about a specific candidate
type Candidate struct {
	Name      string
	ManagerID string
	Meta      map[string]string
}

// The Candidates object is used to manage Candidate instances in a db.Store
// by using the Candidates' names as keys
type Candidates []Candidate

// Get will return the candidate with the matching name.
// The name is not case sensitive.
// A bool is also returned denoting if the candidate exists or not.
func (c Candidates) Get(name string) (Candidate, bool) {
	name = strings.ToLower(name)
	for _, candidate := range c {
		if strings.ToLower(candidate.Name) == name {
			return candidate, true
		}
	}

	return Candidate{}, false
}

// Sort will sort the candidates by their name.
// If ascending is true, names are sorted by alphabetical order.
// If ascending is false, names are sorted by reverse alphabetical order.
func (c Candidates) Sort(ascending bool) {
	if ascending {
		sort.Sort(c)
	} else {
		sort.Sort(sort.Reverse(c))
	}
}

// Len is a method to satisfy sort.Interface
func (c Candidates) Len() int {
	return len(c)
}

// Swap is a method to satisfy sort.Interface
func (c Candidates) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less is a method to satisfy sort.Interface
func (c Candidates) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}
