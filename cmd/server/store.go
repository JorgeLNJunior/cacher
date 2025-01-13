package main

import (
	"sync"
	"time"
)

type InMemoryStore struct {
	data map[string]StoreItem
	mu   sync.Mutex
}

type StoreItem struct {
	Value  string
	Expiry time.Time
}

func (i StoreItem) Expired() bool {
	return time.Now().After(i.Expiry)
}

var oneYear = time.Now().Add(time.Hour * 24 * 365)

// NewInMemoryStore returns a InMemoryStore instance.
func NewInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{
		data: make(map[string]StoreItem),
	}

	go func() {
		for range time.Tick(time.Second * 5) {
			store.mu.Lock()
			for key, item := range store.data {
				if item.Expired() {
					delete(store.data, key)
				}
			}
			store.mu.Unlock()
		}
	}()

	return store
}

// Get returns if the key is stored or not and its value.
func (s *InMemoryStore) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]
	if !ok {
		return "", false
	}

	if item.Expired() {
		delete(s.data, key)
		return "", false
	}

	return item.Value, ok
}

// Set stores a key-value pair in memory.
func (s *InMemoryStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := StoreItem{
		Value:  value,
		Expiry: oneYear,
	}

	s.data[key] = item
}

// Delete removes a key and its value from the storage.
func (s *InMemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

// ExpireAt set an expiration date to an item.
func (s *InMemoryStore) ExpireAt(key string, t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]
	if !ok {
		return
	}

	item.Expiry = t
	s.data[key] = item
}

// Dump returns a copy of all the data in the store.
func (s *InMemoryStore) Dump() map[string]StoreItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	dump := make(map[string]StoreItem)
	for k, v := range s.data {
		dump[k] = v
	}
	return dump
}

func (s *InMemoryStore) Restore(data map[string]StoreItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range data {
		if _, exist := s.data[k]; exist {
			continue
		}
		s.data[k] = v
	}
}
