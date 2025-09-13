package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type RolesStore struct {
	db *pgxpool.Pool
}

func (r *RolesStore) GetByName(ctx context.Context, roleName string) (*Role, error) {
	query := `
		SELECT id, name, description, level FROM roles WHERE name = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := &Role{}
	err := r.db.QueryRow(ctx, query, roleName).
		Scan(&role.ID, &role.Name, &role.Description, &role.Level)
	if err != nil {
		return nil, err
	}

	return role, nil
}
