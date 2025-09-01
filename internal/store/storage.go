package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = 5 * time.Second
)

type Store struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, string) (*Post, error)
		Delete(context.Context, string) error
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		GetByPostID(context.Context, string) ([]Comment, error)
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
		Comments: &CommentsStore{
			db: pool,
		},
	}
}
