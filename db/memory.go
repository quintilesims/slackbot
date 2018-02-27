package db

import (
	"encoding/json"
)

type MemoryStore struct {
	data map[string][]byte
	//Keys() ([]string, error)
	//Read(key string, v interface{}) error
	//Write(key string, v interface{}) error
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: map[string][]byte{},
	}
}

func (m *MemoryStore) Keys() ([]string, error) {
	keys := make([]string, 0, len(m.data))
	for k, _ := range m.data {
		keys = append(keys, k)
	}

	return keys, nil
}

func (m *MemoryStore) Read(key string, v interface{}) error {
	d, ok := m.data[key]
	if !ok {
		return NewMissingEntryError(key)
	}

	return json.Unmarshal(d, &v)
}

func (m *MemoryStore) Write(key string, v interface{}) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}

	m.data[key] = d
	return nil
}
