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
	ErrConflict          = errors.New("resource already exists")
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
		GetByID(context.Context, string) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, string) ([]Comment, error)
	}
	Followers interface {
		Follow(context.Context, string, string) error
		Unfollow(context.Context, string, string) error
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
		Followers: &FollowersStore{
			db: pool,
		},
	}
}
