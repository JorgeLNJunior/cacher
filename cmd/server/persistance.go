package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

type OnDiskStorage struct {
	dataDir      string
	dumpFileName string
	store        *InMemoryStorage
}

// NewOnDiskStorage return a new instance of OnDiskStorage.
func NewOnDiskStorage(store *InMemoryStorage) (*OnDiskStorage, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(userConfigDir, "cacher")
	dumpFileName := "dump"

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return nil, err
	}
	file, err := os.Create(path.Join(dataDir, dumpFileName))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return &OnDiskStorage{
		store:        store,
		dataDir:      dataDir,
		dumpFileName: dumpFileName,
	}, nil
}

// Persist persists the data from memory on disk
func (s OnDiskStorage) Persist(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		dump := s.store.Dump()

		file, err := os.OpenFile(path.Join(s.dataDir, s.dumpFileName), os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(dump); err != nil {
			return err
		}

		return nil
	}
}

// Restore restores the data from disk to memory.
func (s OnDiskStorage) Restore(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		file, err := os.Open(path.Join(s.dataDir, s.dumpFileName))
		if err != nil {
			return err
		}
		defer file.Close()

		dump := make(map[string]StorageItem)
		if err := json.NewDecoder(file).Decode(&dump); err != nil {
			switch {
			case errors.Is(err, io.EOF):
				return nil // file is empty
			default:
				return err
			}
		}
		s.store.Restore(dump)

		return nil
	}
}
