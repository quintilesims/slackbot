package bot

import (
	"fmt"

	"github.com/quintilesims/slack"
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
// Transformations are loaded from the AliasBehavior's store, and then cached for performance.
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

		for name, alias := range aliases {
			if err := alias.Apply(m); err != nil {
				return fmt.Errorf("Alias %s encountered an error: %v", name, err)
			}
		}

		return nil
	}
}

// Invalidate will cause the AliasBehavior to invalidate its cache
func (a *AliasBehavior) Invalidate() {
	a.cache.Clear()
}

func (a *AliasBehavior) load() (models.Aliases, error) {
	if len(a.cache.Keys()) > 0 {
		aliases := models.Aliases{}
		for k, v := range a.cache.Items() {
			aliases[k] = v.(models.Alias)
		}

		return aliases, nil
	}

	aliases := models.Aliases{}
	if err := a.store.Read(db.AliasesKey, &aliases); err != nil {
		return nil, err
	}

	a.cache.Clear()
	for k, v := range aliases {
		a.cache.Add(k, v)
	}

	return aliases, nil
}
