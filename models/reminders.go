package models

import (
	"fmt"
	"time"
)

const StoreKeyReminders = "reminders"

type Reminder struct {
	UserID  string
	Message string
	Time    time.Time
}

func (r Reminder) String() string {
	return fmt.Sprintf("%v %v %v", r.UserID, r.Message, r.Time)
}

type Reminders map[string]Reminder
