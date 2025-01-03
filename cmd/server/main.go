package main

import (
	"log"
	"os"
)

type config struct {
	address string
}

type application struct {
	logger *log.Logger
	config config
}

func main() {
	cfg := config{
		address: ":8595",
	}
	app := &application{
		config: cfg,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	if err := app.Listen(); err != nil {
		app.logger.Fatalf("error listening at %s: %s", app.config.address, err.Error())
	}
}
