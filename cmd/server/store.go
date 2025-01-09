package main

import "sync"

type InMemoryStore struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewInMemoryStore returns a InMemoryStore instance.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]string),
	}
}

// Get returns if the key is stored or not and its value.
func (s *InMemoryStore) Get(key string) (bool, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	return ok, value
}

// Set stores a key-value pair in memory.
func (s *InMemoryStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

// Delete removes a key and its value from the storage.
func (s *InMemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
