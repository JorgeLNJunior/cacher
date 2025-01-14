package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/JorgeLNJunior/cacher/pkg/logger"
)

type config struct {
	address string
}

type application struct {
	config           config
	logger           *logger.Logger
	store            *InMemoryStore
	persistanceStore *OnDiskPersistanceStore
	wg               sync.WaitGroup
}

type loggerArgs map[string]string

func main() {
	store := NewInMemoryStore()
	logger := logger.NewLogger(logger.LevelInfo, os.Stdout)

	persistanceStore, err := NewInDiskPersistanceStore(store)
	if err != nil {
		logger.Error("error creating the on disk persistance store: %s", loggerArgs{"err": err.Error()})
		os.Exit(1)
	}

	logger.Info("restoring the data from disk", nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	if err := persistanceStore.Restore(ctx); err != nil {
		logger.Error("error restoring the data from disk: %s", loggerArgs{"err": err.Error()})
		os.Exit(1)
	}
	cancel()
	logger.Info("the data has been successfully restored", nil)

	cfg := config{
		address: ":8595",
	}

	app := &application{
		config:           cfg,
		logger:           logger,
		store:            store,
		persistanceStore: persistanceStore,
	}

	if err := app.Listen(); err != nil {
		app.logger.Fatal(
			"error listening the server",
			loggerArgs{"addr": app.config.address, "err": err.Error()},
		)
	}
}
