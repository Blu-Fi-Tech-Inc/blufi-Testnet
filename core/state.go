package core

import (
	"errors"
	"fmt"
	"sync"
)

type State struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewState() *State {
	return &State{
		data: make(map[string][]byte),
	}
}

func (s *State) Put(k, v []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(k) == 0 {
		return errors.New("key cannot be empty")
	}

	s.data[string(k)] = v
	return nil
}

func (s *State) Delete(k []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(k) == 0 {
		return errors.New("key cannot be empty")
	}

	delete(s.data, string(k))
	return nil
}

func (s *State) Get(k []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := string(k)
	value, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("given key %s not found", key)
	}

	return value, nil
}
