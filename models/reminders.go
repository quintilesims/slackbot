package models

import (
	"time"
)

const StoreKeyReminders = "reminders"

type Reminder struct {
	UserID  string
	Message string
	Time    time.Time
}

type Reminders map[string]Reminder
