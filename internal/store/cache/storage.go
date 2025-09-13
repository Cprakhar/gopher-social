package cache

import (
	"context"

	"github.com/cprakhar/gopher-social/internal/store"
	"github.com/go-redis/redis/v8"
)

type Store struct {
	Users interface {
		Get(context.Context, string) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewRedisStore(rdb *redis.Client) Store {
	return Store{
		Users: &UserStore{rdb},
	}
}
