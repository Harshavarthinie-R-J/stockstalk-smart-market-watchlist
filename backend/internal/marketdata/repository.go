package marketdata

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveSnapshot(ctx context.Context, q Quote) error {
	receivedTime := q.ReceivedTime

	if receivedTime.IsZero() {
		receivedTime = time.Now().UTC()
	}

	source := q.Source

	if source == "" {
		source = "mock"
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO market_snapshots (
			instrument_id,
			price,
			previous_close,
			open_price,
			high_price,
			low_price,
			volume,
			change_value,
			change_percent,
			week_52_high,
			week_52_low,
			market_status,
			market_timestamp,
			received_timestamp,
			source
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15
		)
		ON CONFLICT (instrument_id, market_timestamp, source)
		DO NOTHING
	`,
		q.InstrumentID,
		q.Price,
		q.PreviousClose,
		q.Open,
		q.High,
		q.Low,
		q.Volume,
		q.Change,
		q.ChangePercent,
		q.Week52High,
		q.Week52Low,
		q.MarketStatus,
		q.MarketTime,
		receivedTime,
		source,
	)

	return err
}

func (r *Repository) SnapshotAtOrBefore(
	ctx context.Context,
	instrumentID string,
	at time.Time,
) (*Quote, error) {

	var q Quote

	err := r.db.QueryRowContext(ctx, `
		SELECT
			instrument_id,
			price,
			previous_close,
			open_price,
			high_price,
			low_price,
			volume,
			change_value,
			change_percent,
			week_52_high,
			week_52_low,
			market_status,
			market_timestamp,
			received_timestamp,
			source
		FROM market_snapshots
		WHERE instrument_id = $1
		  AND market_timestamp <= $2
		ORDER BY market_timestamp DESC, created_at DESC
		LIMIT 1
	`,
		instrumentID,
		at,
	).Scan(
		&q.InstrumentID,
		&q.Price,
		&q.PreviousClose,
		&q.Open,
		&q.High,
		&q.Low,
		&q.Volume,
		&q.Change,
		&q.ChangePercent,
		&q.Week52High,
		&q.Week52Low,
		&q.MarketStatus,
		&q.MarketTime,
		&q.ReceivedTime,
		&q.Source,
	)

	if err != nil {
		return nil, errors.New("snapshot not found")
	}

	return &q, nil
}

func (r *Repository) LatestSnapshot(
	ctx context.Context,
	instrumentID string,
) (*Quote, error) {

	var q Quote

	err := r.db.QueryRowContext(ctx, `
		SELECT
			instrument_id,
			price,
			previous_close,
			open_price,
			high_price,
			low_price,
			volume,
			change_value,
			change_percent,
			week_52_high,
			week_52_low,
			market_status,
			market_timestamp,
			received_timestamp,
			source
		FROM market_snapshots
		WHERE instrument_id = $1
		ORDER BY market_timestamp DESC, created_at DESC
		LIMIT 1
	`,
		instrumentID,
	).Scan(
		&q.InstrumentID,
		&q.Price,
		&q.PreviousClose,
		&q.Open,
		&q.High,
		&q.Low,
		&q.Volume,
		&q.Change,
		&q.ChangePercent,
		&q.Week52High,
		&q.Week52Low,
		&q.MarketStatus,
		&q.MarketTime,
		&q.ReceivedTime,
		&q.Source,
	)

	if err != nil {
		return nil, errors.New("snapshot not found")
	}

	return &q, nil
}

func (r *Repository) AverageVolume(
	ctx context.Context,
	instrumentID string,
	before time.Time,
) (float64, error) {
	var average float64

	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(volume), 0)
		FROM (
			SELECT volume
			FROM market_snapshots
			WHERE instrument_id = $1
			  AND market_timestamp < $2
			  AND volume > 0
			ORDER BY market_timestamp DESC
			LIMIT 20
		) historical
	`, instrumentID, before).Scan(&average)

	if err != nil {
		return 0, err
	}

	return average, nil
}

// SnapshotBefore returns the most recent snapshot strictly before the
// supplied market timestamp.
func (r *Repository) SnapshotBefore(
	ctx context.Context,
	instrumentID string,
	at time.Time,
) (*Quote, error) {
	var q Quote

	err := r.db.QueryRowContext(ctx, `
		SELECT instrument_id, price, previous_close, open_price,
		       high_price, low_price, volume, change_value,
		       change_percent, week_52_high, week_52_low,
		       market_status, market_timestamp,
		       received_timestamp, source
		FROM market_snapshots
		WHERE instrument_id = $1
		  AND market_timestamp < $2
		ORDER BY market_timestamp DESC, created_at DESC
		LIMIT 1
	`,
		instrumentID,
		at,
	).Scan(
		&q.InstrumentID,
		&q.Price,
		&q.PreviousClose,
		&q.Open,
		&q.High,
		&q.Low,
		&q.Volume,
		&q.Change,
		&q.ChangePercent,
		&q.Week52High,
		&q.Week52Low,
		&q.MarketStatus,
		&q.MarketTime,
		&q.ReceivedTime,
		&q.Source,
	)

	if err != nil {
		return nil, errors.New("snapshot not found")
	}

	return &q, nil
}
