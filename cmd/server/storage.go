package main

import "sync"

type InMemoryStorage struct {
	data map[string][]byte
	mu   sync.RWMutex
}

// NewInMemoryStorage returns a InMemoryStorage instance.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string][]byte),
	}
}

// Get returns if the key is stored or not and its value.
func (s *InMemoryStorage) Get(key string) (bool, []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	return ok, value
}

// Set stores a key-value pair in memory.
func (s *InMemoryStorage) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}
