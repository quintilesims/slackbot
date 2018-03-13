package models

// Karma models hold information about a specific karma instance
type Karma struct {
	Upvotes int
	Downvotes int
}

// The Karmas object is used to manage Karma instances in a db.Store
type Karmas map[string]Karma
