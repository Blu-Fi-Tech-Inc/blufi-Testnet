package core

import (
    "fmt"
)

// Store represents a generic key-value store.
type Store interface {
    Get(key string) ([]byte, error)
    Put(key string, value []byte) error
    Delete(key string) error
}

// MemStore is an in-memory key-value store implementation.
type MemStore struct {
    data map[string][]byte
}

// NewMemStore creates a new MemStore.
func NewMemStore() *MemStore {
    return &MemStore{
        data: make(map[string][]byte),
    }
}

// Get retrieves a value by key.
func (m *MemStore) Get(key string) ([]byte, error) {
    value, exists := m.data[key]
    if !exists {
        return nil, fmt.Errorf("key not found: %s", key)
    }
    return value, nil
}

// Put stores a value by key.
func (m *MemStore) Put(key string, value []byte) error {
    m.data[key] = value
    return nil
}

// Delete removes a value by key.
func (m *MemStore) Delete(key string) error {
    delete(m.data, key)
    return nil
}
