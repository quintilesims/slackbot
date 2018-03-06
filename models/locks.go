package models

// StoreKeyLocks is the canonical store key for locks
const StoreKeyLocks = "locks"

// Locks track lock acquisition
// This should be used to read/write locks to/from a db.Store
type Locks map[string]bool
