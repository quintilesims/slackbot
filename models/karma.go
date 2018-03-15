package models

import (
	"fmt"
	"sort"
)

// Karma models hold information about a specific karma instance
type Karma struct {
	Upvotes   int
	Downvotes int
}

// String returns a string representation of a Karma
func (k Karma) String() string {
	return fmt.Sprintf("%d (upvotes: %d, downvotes: %d)",
		k.Upvotes-k.Downvotes,
		k.Upvotes,
		k.Downvotes)
}

// The Karmas object is used to manage Karma instances in a db.Store
type Karmas map[string]Karma

// SortKeys will return a slice of ordered keys.
// If descending is true, keys with the highest karma are returned first.
// If descending is false, keys with the lowest karma are returned first.
func (k Karmas) SortKeys(descending bool) []string {
	sorter := newKarmaSorter(k)
	if descending {
		sort.Sort(sort.Reverse(sorter))
	} else {
		sort.Sort(sorter)
	}

	return sorter.keys
}

type karmaSorter struct {
	karmas Karmas
	keys   []string
}

func newKarmaSorter(karmas Karmas) *karmaSorter {
	keys := make([]string, 0, len(karmas))
	for key := range karmas {
		keys = append(keys, key)
	}

	return &karmaSorter{
		karmas: karmas,
		keys:   keys,
	}
}

func (k *karmaSorter) Len() int {
	return len(k.keys)
}

func (k *karmaSorter) Swap(i, j int) {
	k.keys[i], k.keys[j] = k.keys[j], k.keys[i]
}

func (k *karmaSorter) Less(i, j int) bool {
	karmaI := k.karmas[k.keys[i]]
	karmaJ := k.karmas[k.keys[j]]
	return (karmaI.Upvotes - karmaI.Downvotes) < (karmaJ.Upvotes - karmaJ.Downvotes)
}
