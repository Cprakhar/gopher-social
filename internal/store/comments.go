package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentsStore struct {
	db *pgxpool.Pool
}

type Comment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	AuthorID  string `json:"author_id"`
	Content   string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User   `json:"user"`
}

func (c *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, author_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := c.db.QueryRow(ctx, query, comment.PostID, comment.AuthorID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (c *CommentsStore) GetByPostID(ctx context.Context, postID string) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.author_id, c.content, c.created_at, users.username, users.id FROM comments c
		JOIN users ON c.author_id = users.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := c.db.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		comment.User = User{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt, &comment.User.Username, &comment.User.ID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
