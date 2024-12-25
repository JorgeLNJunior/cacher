package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"testing"
)

func TestSet(t *testing.T) {
	storage := NewInMemoryStorage()
	expectedKeys := 100

	for i := 0; i < expectedKeys; i++ {
		storage.Set(string(randomBytes()), randomBytes())
	}

	insertedKeys := len(storage.data)

	if insertedKeys < expectedKeys {
		t.Fatalf("expected %d inserted keys but got %d", expectedKeys, insertedKeys)
	}
}

func TestGet(t *testing.T) {
	storage := NewInMemoryStorage()

	key := "foo"
	value := []byte("bar")

	storage.Set(key, value)

	if !bytes.Equal(storage.data[key], value) {
		t.Fatal("inserted value differs from retrieved value")
	}
}

func randomBytes() []byte {
	data := make([]byte, 16)
	_, _ = rand.Read(data)
	base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
	return data
}
