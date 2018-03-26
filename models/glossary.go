package models

import "sort"

// The Glossary object is used to manage a Glossary in a db.Store
type Glossary map[string]string

// SortKeys will return a slice of alphabetically ordered keys.
// If ascending is true, keys are returned in alphabetical order.
// If ascending is false, keys are returned in reverse alphabetical order.
func (g Glossary) SortKeys(ascending bool) []string {
	keys := make([]string, 0, len(g))
	for key := range g {
		keys = append(keys, key)
	}

	if ascending {
		sort.Sort(sort.StringSlice(keys))
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}

	return keys
}
