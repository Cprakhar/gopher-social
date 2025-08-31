package main

import (
	"context"
	"log"
	"time"

	"github.com/cprakhar/gopher-social/internal/config"
	"github.com/cprakhar/gopher-social/internal/handler"
	"github.com/cprakhar/gopher-social/internal/db"
	"github.com/cprakhar/gopher-social/internal/store"
)

type application struct {
	config  config.Config
	store   store.Store
	handler handler.Handler
}

func main() {

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := db.New(ctx,
		cfg.DB.Addr,
		cfg.DB.MaxConns,
		cfg.DB.MinConns,
		cfg.DB.MaxIdleTime,
		cfg.DB.MaxConnLifetime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	log.Println("database connection pool established")

	store := store.NewStore(db)

	app := &application{
		config: cfg,
		store:  store,
		handler: handler.Handler{
			Cfg: cfg,
		},
	}

	mux := app.mount()

	log.Printf("server is running on %s", cfg.Addr)
	log.Fatal(app.run(mux))
}
