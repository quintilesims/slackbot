package db

import "fmt"

type MissingEntryError error

func NewMissingEntryError(key string) MissingEntryError {
	return MissingEntryError(fmt.Errorf("No entry for key '%s'", key))
}
