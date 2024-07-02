package core

import (
    "fmt"
    "sync"
)

type Storage interface {
    Put(*Block) error
    Get(hash []byte) (*Block, error)
}

type MemoryStore struct {
    mu    sync.RWMutex
    store map[string]*Block
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        store: make(map[string]*Block),
    }
}

func (s *MemoryStore) Put(b *Block) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    hash := b.Hash(BlockHasher{}).String()
    s.store[hash] = b
    return nil
}

func (s *MemoryStore) Get(hash []byte) (*Block, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    b, ok := s.store[string(hash)]
    if !ok {
        return nil, fmt.Errorf("block not found")
    }
    return b, nil
}
