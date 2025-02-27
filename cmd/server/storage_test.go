package main

import (
	"crypto/rand"
	"encoding/base32"
	"testing"
)

func TestSet(t *testing.T) {
	storage := NewInMemoryStorage()
	expectedKeys := 100

	for range expectedKeys {
		storage.Set(string(randomString()), randomString())
	}

	insertedKeys := len(storage.data)

	if insertedKeys < expectedKeys {
		t.Fatalf("expected %d inserted keys but got %d", expectedKeys, insertedKeys)
	}
}

func TestGet(t *testing.T) {
	storage := NewInMemoryStorage()

	key := "foo"
	value := "bar"

	storage.Set(key, value)

	if storage.data[key].Value != value {
		t.Fatal("inserted value differs from retrieved value")
	}
}

func TestDelete(t *testing.T) {
	storage := NewInMemoryStorage()
	key := "foo"

	storage.Set(key, randomString())
	storage.Delete(key)

	if _, ok := storage.Get(key); ok {
		t.Fatal("the key has not been deleted")
	}
}

func BenchmarkSet(b *testing.B) {
	storage := NewInMemoryStorage()
	for b.Loop() {
		storage.Set(randomString(), randomString())
	}
}

func BenchmarkGet(b *testing.B) {
	storage := NewInMemoryStorage()

	key := randomString()
	storage.Set(key, randomString())

	for b.Loop() {
		storage.Get(key)
	}
}

func randomString() string {
	data := make([]byte, 16)
	_, _ = rand.Read(data)
	base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
	return string(data)
}
