package db

import "errors"

type MissingEntryError error

func NewMissingEntryError(key string) MissingEntryError {
	return MissingEntryError(errors.New("No entry for key '%s'", key))
}
