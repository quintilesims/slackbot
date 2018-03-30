package models

import (
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
