package main

import (
	"log"
	"os"
	"sync"
)

type config struct {
	address string
}

type application struct {
	config config
	logger *log.Logger
	store  *InMemoryStore
	wg     sync.WaitGroup
}

func main() {
	cfg := config{
		address: ":8595",
	}
	app := &application{
		config: cfg,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		store:  NewInMemoryStore(),
	}

	if err := app.Listen(); err != nil {
		app.logger.Printf("error listening at %s: %s\n", app.config.address, err.Error())
	}
}
