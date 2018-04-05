package models

import (
	"sort"
	"strings"
	"time"
)

// Interview models hold information about a single interview instance
type Interview struct {
	Candidate      string
	InterviewerIDs []string
	Time           time.Time
}

func (i Interview) Equals(other Interview) bool {
	if strings.ToLower(i.Candidate) != strings.ToLower(other.Candidate) {
		return false
	}

	return i.Time == other.Time
}

// The Interviews object is used to manage Interview instances in a db.Store
type Interviews []Interview

// Sort will sort the interviews by their time.
// If ascending is true, names are sorted in chronological order.
// If ascending is false, names are sorted by reverse chronological order.
func (n Interviews) Sort(ascending bool) {
	if ascending {
		sort.Sort(n)
	} else {
		sort.Sort(sort.Reverse(n))
	}
}

// Len is a method to satisfy sort.Interface
func (n Interviews) Len() int {
	return len(n)
}

// Swap is a method to satisfy sort.Interface
func (n Interviews) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// Less is a method to satisfy sort.Interface
func (n Interviews) Less(i, j int) bool {
	return n[i].Time.Before(n[j].Time)
}
