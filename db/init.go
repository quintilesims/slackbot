package db

import "github.com/quintilesims/slackbot/models"

// Init will initialize the table entries for the specified store
func Init(store Store) error {
	initFunc := func(key string, v interface{}) error {
		if err := store.Read(key, &v); err != nil {
			if _, ok := err.(MissingEntryError); ok {
				return store.Write(key, v)
			}

			return err
		}

		return nil
	}

	if err := initFunc(AliasesKey, models.Aliases{}); err != nil {
		return err
	}

	if err := initFunc(CandidatesKey, models.Candidates{}); err != nil {
		return err
	}

	if err := initFunc(GlossaryKey, models.Glossary{}); err != nil {
		return err
	}

	if err := initFunc(InterviewsKey, models.Interviews{}); err != nil {
		return err
	}

	if err := initFunc(KarmasKey, models.Karmas{}); err != nil {
		return err
	}

	return nil
}
