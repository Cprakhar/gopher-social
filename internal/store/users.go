package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersStore struct {
	db *pgxpool.Pool
}

func (p *UsersStore) Create(ctx context.Context) error {
	return nil
}