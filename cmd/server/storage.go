package main

import (
	"maps"
	"sync"
	"time"
)

var oneYear = time.Now().Add(time.Hour * 24 * 365)

type Storage interface {
	Restore(data map[string]StorageItem)
	Dump() map[string]StorageItem
}

type InMemoryStorage struct {
	data map[string]StorageItem
	mu   sync.Mutex
}

type StorageItem struct {
	Value  string
	Expiry time.Time
}

// Expired returns whether an item is expired.
func (i StorageItem) Expired() bool {
	return time.Now().After(i.Expiry)
}

// NewInMemoryStorage returns a InMemoryStorage instance.
func NewInMemoryStorage() *InMemoryStorage {
	store := &InMemoryStorage{
		data: make(map[string]StorageItem),
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
func (s *InMemoryStorage) Get(key string) (string, bool) {
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

// Set stores a key-value pair into the store.
func (s *InMemoryStorage) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := StorageItem{
		Value:  value,
		Expiry: oneYear,
	}

	s.data[key] = item
}

// Delete removes a key and its value from the storage.
func (s *InMemoryStorage) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

// ExpireAt sets the expiration date of an item.
func (s *InMemoryStorage) ExpireAt(key string, t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, found := s.data[key]
	if !found {
		return
	}

	item.Expiry = t
	s.data[key] = item
}

// Dump returns a copy of all data in the storage.
func (s *InMemoryStorage) Dump() map[string]StorageItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	dump := make(map[string]StorageItem, len(s.data))
	maps.Copy(dump, s.data)

	return dump
}

// Restore copies data into the storage.
func (s *InMemoryStorage) Restore(data map[string]StorageItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range data {
		if _, found := s.data[k]; !found {
			s.data[k] = v
		}
	}
}
