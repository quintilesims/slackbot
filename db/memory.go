package db

import (
	"encoding/json"
)

// MemoryStore reads and writes data to memory
type MemoryStore struct {
	data map[string][]byte
}

// NewMemoryStore creates a new MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: map[string][]byte{},
	}
}

// Keys lists all of the keys in the store
func (m *MemoryStore) Keys() ([]string, error) {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}

	return keys, nil
}

// Read will populate v with the entry at the specified key
func (m *MemoryStore) Read(key string, v interface{}) error {
	d, ok := m.data[key]
	if !ok {
		return NewMissingEntryError(key)
	}

	return json.Unmarshal(d, &v)
}

// Write will write v at the specified key
func (m *MemoryStore) Write(key string, v interface{}) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}

	m.data[key] = d
	return nil
}
