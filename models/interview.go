package models

import (
	"time"
)

// Interview hold information about an interview
type Interview struct {
	ManagerID   string
	ManagerName string
	Interviewee string
	Date        time.Time
}

// Interviews track Interview objects using the interview's ID as the key
type Interviews map[string]Interview
