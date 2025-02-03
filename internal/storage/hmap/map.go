package hmap

import (
	"sync"
)

// MapStore is an in-memory implementation of Storage using sync.Map
type MapStore struct {
	store sync.Map
}


// NewMapStore creates a new MapStore instance
func NewMapStore() *MapStore {
	return &MapStore{}
}

// Get retrieves a value from the map
func (m *MapStore) Get(key string) (interface{}, bool) {
	return m.store.Load(key)
}

// Put stores a value in the map
func (m *MapStore) Put(key string, value interface{}) {
	m.store.Store(key, value)
}

// Pop deletes a key from the map
func (m *MapStore) Pop(key string) {
	m.store.Delete(key)
}
