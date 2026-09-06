package user

import (
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, name, email, passwordHash string) (*User, error) {
	var u User

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, created_at, updated_at
	`, name, email, passwordHash).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("email already registered")
	}

	return &u, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, string, error) {
	var u User
	var passwordHash string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&passwordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, "", errors.New("user not found")
	}

	return &u, passwordHash, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	var u User

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &u, nil
}
