package models

import (
	"time"
)

// Reminder hold information about a reminder
type Reminder struct {
	UserID   string
	UserName string
	Message  string
	Time     time.Time
}

// Reminders track Reminder objects using the reminder's ID as the key
type Reminders map[string]Reminder
