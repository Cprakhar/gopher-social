package main

import (
	"log"

	"github.com/cprakhar/gopher-social/handler"
	"github.com/cprakhar/gopher-social/internal/env"
	"github.com/cprakhar/gopher-social/internal/store"
)

type config struct {
	addr string
}

type application struct {
	config config
	store store.Store
	handler *handler.Handler
}

func main() {

	config := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStore(nil)

	app := &application{
		config: config,
		store:  store,
	}

	mux := app.mount()

	log.Printf("server is running on %s", config.addr)
	log.Fatal(app.run(mux))
}
