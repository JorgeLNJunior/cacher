package main

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

type OnDiskPersistanceStore struct {
	dataDir      string
	dumpFileName string
	store        *InMemoryStore
}

// NewInDiskPersistanceStore return a new instance of OnDiskPersistanceStore.
func NewInDiskPersistanceStore(store *InMemoryStore) (*OnDiskPersistanceStore, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(userConfigDir, "cacher")
	dumpFileName := "dump"

	return &OnDiskPersistanceStore{
		store:        store,
		dataDir:      dataDir,
		dumpFileName: dumpFileName,
	}, nil
}

// Persist persists the data from memory on disk
func (s OnDiskPersistanceStore) Persist(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		dump := s.store.Dump()

		if err := os.MkdirAll(s.dataDir, os.ModePerm); err != nil {
			return err
		}

		file, err := os.Create(path.Join(s.dataDir, s.dumpFileName))
		if err != nil {
			return err
		}

		data, err := json.Marshal(dump)
		if err != nil {
			return err
		}

		if _, err := file.Write(data); err != nil {
			return err
		}

		return nil
	}
}

// Restore restores the data from disk to memory.
func (s OnDiskPersistanceStore) Restore(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		file, err := os.Open(path.Join(s.dataDir, s.dumpFileName))
		if err != nil {
			return err
		}

		dump := make(map[string]StoreItem)
		if err := json.NewDecoder(file).Decode(&dump); err != nil {
			return err
		}
		s.store.Restore(dump)

		return nil
	}
}
