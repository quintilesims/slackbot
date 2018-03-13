package models

import (
	"time"
)

// Interview models hold information about a single interview instance
type Interview struct {
	ManagerID   string
	ManagerName string
	Interviewee string
	Date        time.Time
}

// The Interviews object is used to manage Interview instances in a db.Store
type Interviews map[string]Interview
