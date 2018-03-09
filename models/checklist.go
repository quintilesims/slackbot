package models

// ChecklistItem holds information about a single item in a checklist
type ChecklistItem struct {
	ID        string
	Text      string
	Source    string
	IsChecked bool
}

// Checklist stores a slice of CheckListItem objects
type Checklist []ChecklistItem

// Checklists track Checklist objects using the checklist owner's user ID as the key
type Checklists map[string]Checklist
