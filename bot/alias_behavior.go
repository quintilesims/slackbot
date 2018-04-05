package bot

import (
	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	cache "github.com/zpatrick/go-cache"
)

// AliasBehavior will apply alias transformations to slack message events
type AliasBehavior struct {
	store db.Store
	cache *cache.Cache
}

// NewAliasBehavior will create a new instance of an AliasBehavior object
func NewAliasBehavior(store db.Store) *AliasBehavior {
	return &AliasBehavior{
		store: store,
		cache: cache.New(),
	}
}

// Behavior will return the AliasBehavior's Behavior function.
// This function will apply alias transformations to slack message events.
// Aliases are loaded from the AliasBehavior's store, and then cached for performance.
// The Invalidate() function should be called to invalid the AliasBehavior's cache.
func (a *AliasBehavior) Behavior() Behavior {
	return func(e slack.RTMEvent) error {
		m, ok := e.Data.(*slack.MessageEvent)
		if !ok {
			return nil
		}

		aliases, err := a.load()
		if err != nil {
			return err
		}

		return aliases.Apply(m)
	}
}

// Invalidate will cause the AliasBehavior to invalidate its cache
func (a *AliasBehavior) Invalidate() {
	a.cache.Clear()
}

func (a *AliasBehavior) load() (models.Aliases, error) {
	if v, ok := a.cache.Getf("key"); ok {
		return v.(models.Aliases), nil
	}

	aliases := models.Aliases{}
	if err := a.store.Read(db.AliasesKey, &aliases); err != nil {
		return nil, err
	}

	a.cache.Set("key", aliases)
	return aliases, nil
}
