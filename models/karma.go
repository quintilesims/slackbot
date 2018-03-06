package models

// StoreKeyKarma is the canonical store key for karma
const StoreKeyKarma = "karma"

// Karma tracks karma scores
// This should be used to read/write karma to/from a db.Store
type Karma map[string]int
