package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type FollowersStore struct {
	db *pgxpool.Pool
}

type Follower struct {
	UserID      string    `json:"user_id"`
	FollowingID string    `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (u *FollowersStore) Follow(ctx context.Context, followingID, id string) error {
	query := `
		INSERT INTO followers (user_id, following_id)
		VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := u.db.Exec(ctx, query, id, followingID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
		return err	
	}
	return nil
}

func (u *FollowersStore) Unfollow(ctx context.Context, followingID, id string) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND following_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := u.db.Exec(ctx, query, id, followingID)
	return err
}
