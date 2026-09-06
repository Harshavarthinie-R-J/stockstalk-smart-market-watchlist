package db

import (
	"context"
	"database/sql"

	"stockstalk/internal/marketdata"
)

func SeedInstruments(ctx context.Context, conn *sql.DB, instruments []marketdata.Instrument) error {
	for _, item := range instruments {
		_, err := conn.ExecContext(ctx, `
			INSERT INTO instruments (
				id,
				symbol,
				name,
				exchange,
				segment,
				isin,
				currency,
				sector,
				industry,
				is_active
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			ON CONFLICT (id) DO UPDATE SET
				symbol = EXCLUDED.symbol,
				name = EXCLUDED.name,
				exchange = EXCLUDED.exchange,
				segment = EXCLUDED.segment,
				isin = EXCLUDED.isin,
				currency = EXCLUDED.currency,
				sector = EXCLUDED.sector,
				industry = EXCLUDED.industry,
				is_active = EXCLUDED.is_active
		`,
			item.ID,
			item.Symbol,
			item.Name,
			item.Exchange,
			item.Segment,
			item.ISIN,
			item.Currency,
			item.Sector,
			item.Industry,
			true,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
