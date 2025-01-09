package main

import (
	"crypto/rand"
	"encoding/base32"
	"testing"
)

func TestSet(t *testing.T) {
	storage := NewInMemoryStore()
	expectedKeys := 100

	for i := 0; i < expectedKeys; i++ {
		storage.Set(string(randomString()), randomString())
	}

	insertedKeys := len(storage.data)

	if insertedKeys < expectedKeys {
		t.Fatalf("expected %d inserted keys but got %d", expectedKeys, insertedKeys)
	}
}

func TestGet(t *testing.T) {
	storage := NewInMemoryStore()

	key := "foo"
	value := "bar"

	storage.Set(key, value)

	if storage.data[key] != value {
		t.Fatal("inserted value differs from retrieved value")
	}
}

func TestDelete(t *testing.T) {
	storage := NewInMemoryStore()
	key := "foo"

	storage.Set(key, randomString())
	storage.Delete(key)

	if _, ok := storage.Get(key); ok {
		t.Fatal("the key has not been deleted")
	}
}

func randomString() string {
	data := make([]byte, 16)
	_, _ = rand.Read(data)
	base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
	return string(data)
}
