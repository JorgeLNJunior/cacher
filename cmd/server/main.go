package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

type config struct {
	address string
}

type application struct {
	config           config
	logger           *log.Logger
	store            *InMemoryStore
	persistanceStore *OnDiskPersistanceStore
	wg               sync.WaitGroup
}

func main() {
	store := NewInMemoryStore()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	persistanceStore, err := NewInDiskPersistanceStore(store)
	if err != nil {
		logger.Printf("error creating the on disk persistance store: %s", err)
		os.Exit(1)
	}

	logger.Println("restoring the data from disk")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	if err := persistanceStore.Restore(ctx); err != nil {
		logger.Printf("error restoring the data from disk: %s", err)
		os.Exit(1)
	}
	cancel()
	logger.Println("the data has been successfully restored")

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
		app.logger.Printf("error listening at %s: %s\n", app.config.address, err.Error())
	}
}
