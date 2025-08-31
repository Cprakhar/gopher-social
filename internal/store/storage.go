package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Posts interface {
		Create(context.Context) error
	}
	Users interface {
		Create(context.Context) error
	}
}

func NewStore(pool *pgxpool.Pool) Store {
	return Store{
		Posts: &PostsStore{
			db: pool,
		},
		Users: &UsersStore{
			db: pool,
		},
	}
}
