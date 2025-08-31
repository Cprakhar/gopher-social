package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostsStore struct {
	db *pgxpool.Pool
}

func (p *PostsStore) Create(ctx context.Context) error {
	return nil
}