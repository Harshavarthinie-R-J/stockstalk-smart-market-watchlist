package checkpoint

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

func (r *Repository) Create(
	ctx context.Context,
	userID,
	watchlistID string,
) (*Checkpoint, error) {

	var c Checkpoint

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO checkpoints (
			user_id,
			watchlist_id
		)
		VALUES ($1, $2)
		RETURNING id, user_id, watchlist_id, created_at
	`,
		userID,
		watchlistID,
	).Scan(
		&c.ID,
		&c.UserID,
		&c.WatchlistID,
		&c.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *Repository) Latest(
	ctx context.Context,
	userID,
	watchlistID string,
) (*Checkpoint, error) {

	var c Checkpoint

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			user_id,
			watchlist_id,
			created_at
		FROM checkpoints
		WHERE user_id = $1
		  AND watchlist_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`,
		userID,
		watchlistID,
	).Scan(
		&c.ID,
		&c.UserID,
		&c.WatchlistID,
		&c.CreatedAt,
	)

	if err != nil {
		return nil, errors.New("checkpoint not found")
	}

	return &c, nil
}
