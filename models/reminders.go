package models

import (
	"fmt"
	"time"
)

const StoreKeyReminders = "reminders"

// Reminders hold information about a reminder
type Reminder struct {
	UserID   string
	UserName string
	Message  string
	Time     time.Time
}

// String returns the string reprentation of the reminder
func (r Reminder) String() string {
	return fmt.Sprintf("%v %v %v %v", r.UserID, r.UserName, r.Message, r.Time)
}

// Reminders track Reminder objects using the reminder's ID as the key
// This should be used to read/write reminders to/from a db.Store
type Reminders map[string]Reminder
