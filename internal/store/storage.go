package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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
		GetUserFeed(context.Context, string, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, pgx.Tx, *User) error
		Delete(context.Context, string) error
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		GetByID(context.Context, string) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Activate(context.Context, string) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, string) ([]Comment, error)
	}
	Followers interface {
		Follow(context.Context, string, string) error
		Unfollow(context.Context, string, string) error
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

func NewStore(db *pgxpool.Pool) Store {
	return Store{
		Posts:     &PostsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowersStore{db},
		Roles:     &RolesStore{db},
	}
}

func withTx(db *pgxpool.Pool, ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
