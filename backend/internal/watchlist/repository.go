package watchlist

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

func (r *Repository) Create(ctx context.Context, userID, name string) (*Watchlist, error) {
	var w Watchlist

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO watchlists (user_id, name)
		VALUES ($1, $2)
		RETURNING id, user_id, name, created_at, updated_at
	`, userID, name).Scan(
		&w.ID,
		&w.UserID,
		&w.Name,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	w.Stocks = []WatchItem{}

	return &w, nil
}

func (r *Repository) GetByID(ctx context.Context, userID, id string) (*Watchlist, error) {
	var w Watchlist

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, name, created_at, updated_at
		FROM watchlists
		WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(
		&w.ID,
		&w.UserID,
		&w.Name,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("watchlist not found")
	}

	stocks, err := r.getStocks(ctx, w.ID)
	if err != nil {
		return nil, err
	}

	w.Stocks = stocks

	return &w, nil
}

func (r *Repository) List(ctx context.Context, userID string) ([]Watchlist, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, created_at, updated_at
		FROM watchlists
		WHERE user_id = $1
		ORDER BY created_at
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]Watchlist, 0)

	for rows.Next() {
		var w Watchlist

		if err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.Name,
			&w.CreatedAt,
			&w.UpdatedAt,
		); err != nil {
			return nil, err
		}

		stocks, err := r.getStocks(ctx, w.ID)
		if err != nil {
			return nil, err
		}

		w.Stocks = stocks

		result = append(result, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) AddStock(ctx context.Context, userID, watchlistID, instrumentID string) error {
	var exists bool

	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM watchlists
			WHERE id = $1 AND user_id = $2
		)
	`, watchlistID, userID).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("watchlist not found")
	}

	var position int

	err = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(position) + 1, 0)
		FROM watchlist_stocks
		WHERE watchlist_id = $1
	`, watchlistID).Scan(&position)

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO watchlist_stocks (
			watchlist_id,
			instrument_id,
			position
		)
		VALUES ($1, $2, $3)
	`, watchlistID, instrumentID, position)

	if err != nil {
		return errors.New("stock already in watchlist")
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE watchlists
		SET updated_at = NOW()
		WHERE id = $1
	`, watchlistID)

	return err
}

func (r *Repository) RemoveStock(ctx context.Context, userID, watchlistID, instrumentID string) error {
	var exists bool

	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM watchlists
			WHERE id = $1 AND user_id = $2
		)
	`, watchlistID, userID).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("watchlist not found")
	}

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM watchlist_stocks
		WHERE watchlist_id = $1
		AND instrument_id = $2
	`, watchlistID, instrumentID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("stock not found in watchlist")
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE watchlists
		SET updated_at = NOW()
		WHERE id = $1
	`, watchlistID)

	return err
}

func (r *Repository) getStocks(ctx context.Context, watchlistID string) ([]WatchItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT instrument_id, position, added_at
		FROM watchlist_stocks
		WHERE watchlist_id = $1
		ORDER BY position
	`, watchlistID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stocks := make([]WatchItem, 0)

	for rows.Next() {
		var item WatchItem

		if err := rows.Scan(
			&item.InstrumentID,
			&item.Position,
			&item.AddedAt,
		); err != nil {
			return nil, err
		}

		stocks = append(stocks, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}

// ListAll returns every watchlist in the system.
// This is used by the background market processor.
func (r *Repository) ListAll(ctx context.Context) ([]Watchlist, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, created_at, updated_at
		FROM watchlists
		ORDER BY created_at
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]Watchlist, 0)

	for rows.Next() {
		var w Watchlist

		if err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.Name,
			&w.CreatedAt,
			&w.UpdatedAt,
		); err != nil {
			return nil, err
		}

		stocks, err := r.getStocks(ctx, w.ID)
		if err != nil {
			return nil, err
		}

		w.Stocks = stocks

		result = append(result, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
