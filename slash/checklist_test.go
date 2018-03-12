package slash

import (
	"testing"

	"github.com/quintilesims/slackbot/models"
	"github.com/stretchr/testify/assert"
)

func TestChecklistShow(t *testing.T) {
	checklists := models.Checklists{
		"uid": models.Checklist{
			models.ChecklistItem{ID: "item1"},
			models.ChecklistItem{ID: "item2"},
		},
	}

	store := newMemoryStore(t)
	if err := store.Write(models.StoreKeyChecklists, checklists); err != nil {
		t.Fatal(err)
	}

	result, err := checklistShow(store, "uid")
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result.Attachments, 2)
	assert.Equal(t, result.Attachments[0].CallbackID, "uid.item1")
	assert.Equal(t, result.Attachments[1].CallbackID, "uid.item2")
}

func TestChecklistAdd(t *testing.T) {
	store := newMemoryStore(t)
	if _, err := checklistAdd(store, "uid", "some message"); err != nil {
		t.Fatal(err)
	}

	checklists := models.Checklists{}
	if err := store.Read(models.StoreKeyChecklists, &checklists); err != nil {
		t.Fatal(err)
	}

	result := checklists["uid"]
	assert.Len(t, result, 1)
	assert.NotNil(t, result[0].ID)
	assert.Equal(t, "some message", result[0].Text)
	assert.False(t, result[0].IsChecked)
}

func TestChecklistCallback(t *testing.T) {
	// test actions "check", "uncheck", and "delete"
	t.Skip("TODO")
}
