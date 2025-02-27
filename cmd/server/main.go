package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/JorgeLNJunior/cacher/pkg/logger"
)

type config struct {
	address string
	persist bool
}

type application struct {
	config             config
	logger             *logger.Logger
	storage            *InMemoryStorage
	persistanceStorage *OnDiskStorage
	connectionGroup    sync.WaitGroup
}

type loggerArgs map[string]string

func main() {
	cfg := config{}

	flag.StringVar(&cfg.address, "address", ":8595", "address tcp server will listen")
	flag.BoolVar(&cfg.persist, "persist", false, "persist data on disk or not")
	flag.Parse()

	storage := NewInMemoryStorage()
	logger := logger.NewLogger(logger.LevelInfo, os.Stdout)

	persistanceStorage, err := NewOnDiskStorage()
	if err != nil {
		logger.Fatal("error creating the on disk persistance store: %s", loggerArgs{"err": err.Error()})
	}

	app := &application{
		config:             cfg,
		logger:             logger,
		storage:            storage,
		persistanceStorage: persistanceStorage,
	}

	if app.config.persist {
		logger.Info("restoring the data from disk", nil)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		if err := persistanceStorage.Restore(ctx, app.storage); err != nil {
			logger.Fatal("error restoring the data from disk: %s", loggerArgs{"err": err.Error()})
		}
		cancel()

		logger.Info("the data has been successfully restored", nil)
	}

	if err := app.Listen(); err != nil {
		app.logger.Fatal(
			"error listening the server",
			loggerArgs{"addr": app.config.address, "err": err.Error()},
		)
	}
}
