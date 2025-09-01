package main

import (
	"context"
	"time"

	"github.com/cprakhar/gopher-social/internal/db"
	"github.com/cprakhar/gopher-social/internal/env"
	"github.com/cprakhar/gopher-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gopher-social?sslmode=disable")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := db.New(ctx, addr, 3, 3, time.Minute, time.Minute)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	store := store.NewStore(conn)

	db.Seed(store)
}
