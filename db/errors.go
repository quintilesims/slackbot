package db

import "github.com/quintilesims/layer0/common/errors"

type MissingEntryError error

func NewMissingEntryError(key string) MissingEntryError {
	return MissingEntryError(errors.Newf("No entry for key '%s'", key))
}
