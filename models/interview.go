package models

import (
	"time"
)

// Interview models hold information about a single interview instance
type Interview struct {
	Manager     User
	Interviewee string
	Date        time.Time
}

// The Interviews object is used to manage Interview instances in a db.Store
type Interviews []Interview
