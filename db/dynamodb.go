package db

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type entry struct {
	Key   string
	Value string
}

type DynamoDBStore struct {
	table dynamo.Table
}

func NewDynamoDBStore(session *session.Session, table string) *DynamoDBStore {
	return &DynamoDBStore{
		table: dynamo.New(session).Table(table),
	}
}

func (d *DynamoDBStore) Keys() ([]string, error) {
	entries := []entry{}
	if err := d.table.Scan().
		Consistent(false).
		All(&entries); err != nil {
		return nil, err
	}

	keys := make([]string, len(entries))
	for i, entry := range entries {
		keys[i] = entry.Key
	}

	return keys, nil
}

func (d *DynamoDBStore) Read(key string, v interface{}) error {
	var e entry
	if err := d.table.Get("Key", key).Consistent(true).One(&e); err != nil {
		if err.Error() == "dynamo: no item found" {
			return NewMissingEntryError(key)
		}

		return err
	}

	return json.Unmarshal([]byte(e.Value), &v)
}

func (d *DynamoDBStore) Write(key string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	e := entry{Key: key, Value: string(b)}
	return d.table.Put(e).Run()
}
