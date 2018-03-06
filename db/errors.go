package db

import "fmt"

// MissingEntryError occurs when a Read operation runs with a key that has no corresponding entry
type MissingEntryError error

// NewMissingEntryError creates a new MissingEntryError object
func NewMissingEntryError(key string) MissingEntryError {
	return MissingEntryError(fmt.Errorf("No entry for key '%s'", key))
}
