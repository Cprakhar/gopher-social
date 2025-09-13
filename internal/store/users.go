package store

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	Activated bool      `json:"activated"`
	Role      Role      `json:"role"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

type UsersStore struct {
	db *pgxpool.Pool
}

func (u *UsersStore) Create(ctx context.Context, tx pgx.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4))
		RETURNING id, created_at
	`

	role := user.Role.Name
	if role == "" {
		role = "user"
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := tx.QueryRow(ctx, query, user.Username, user.Email, user.Password.hash, role).
		Scan(&user.ID, &user.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (u *UsersStore) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT users.id, username, email, password, created_at, roles.*
		FROM users
		JOIN roles ON users.role_id = roles.id
		WHERE users.id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := u.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password.hash, &user.CreatedAt, &user.Role.ID, &user.Role.Name, &user.Role.Description, &user.Role.Level)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE email = $1 AND activated = true
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := u.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password.hash, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *UsersStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(u.db, ctx, func(tx pgx.Tx) error {
		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}

		// create the user invitation
		err := u.createUserInvitation(ctx, tx, user.ID, token, invitationExp)
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *UsersStore) Activate(ctx context.Context, token string) error {
	return withTx(u.db, ctx, func(tx pgx.Tx) error {
		user, err := u.getByToken(ctx, tx, token)
		if err != nil {
			return err
		}
		user.Activated = true
		if err := u.update(ctx, tx, user); err != nil {
			return err
		}

		if err := u.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (u *UsersStore) Delete(ctx context.Context, id string) error {
	return withTx(u.db, ctx, func(tx pgx.Tx) error {
		if err := u.delete(ctx, tx, id); err != nil {
			return err
		}

		if err := u.deleteUserInvitation(ctx, tx, id); err != nil {
			return err
		}

		return nil
	})
}

func (u *UsersStore) delete(ctx context.Context, tx pgx.Tx, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	cmdTag, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (u *UsersStore) deleteUserInvitation(ctx context.Context, tx pgx.Tx, userID string) error {
	query := `
		DELETE FROM user_invitations
		WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	cmdTag, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (u *UsersStore) update(ctx context.Context, tx pgx.Tx, user *User) error {
	query := `
		UPDATE users
		SET activated = $1
		WHERE id = $2
		RETURNING id
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRow(ctx, query, user.Activated, user.ID).Scan(&user.ID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (u *UsersStore) createUserInvitation(ctx context.Context, tx pgx.Tx, userID, token string, exp time.Duration) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.Exec(ctx, query, []byte(token), userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (u *UsersStore) getByToken(ctx context.Context, tx pgx.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at
		FROM users u
		JOIN user_invitations ui ON ui.user_id = u.id
		WHERE ui.token = $1 AND ui.expires_at > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := tx.QueryRow(ctx, query, hashToken, time.Now()).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
