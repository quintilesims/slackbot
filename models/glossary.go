package models

import "sort"

// The Glossary object is used to manage a Glossary in a db.Store
type Glossary map[string]string

// SortKeys will return a slice of ordered keys.
func (g Glossary) SortKeyAlphabetical() []string {
	m := make([]string, len(g))
	i := 0
	for k := range g {
		m[i] = k
		i++
	}
	sort.Strings(m)
	return m
}
