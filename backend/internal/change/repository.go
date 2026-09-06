package change

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	instrumentID string,
	c MarketChange,
) error {

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO market_changes (
			instrument_id,
			change_type,
			severity,
			previous_value,
			current_value,
			change_percent,
			current_volume,
			average_volume,
			volume_ratio,
			volume_anomaly,
			price_anomaly,
			near_52w_high,
			near_52w_low,
			signals,
			context,
			confidence,
			attention_score,
			context_type,
			attribution,
			attribution_confidence,
			attribution_source,
			detected_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
			$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22
		)
		ON CONFLICT (
			instrument_id,
			change_type,
			previous_value,
			current_value,
			change_percent
		)
		DO UPDATE SET
			current_volume = EXCLUDED.current_volume,
			average_volume = EXCLUDED.average_volume,
			volume_ratio = EXCLUDED.volume_ratio,
			volume_anomaly = EXCLUDED.volume_anomaly,
			price_anomaly = EXCLUDED.price_anomaly,
			near_52w_high = EXCLUDED.near_52w_high,
			near_52w_low = EXCLUDED.near_52w_low,
			signals = EXCLUDED.signals,
			context = EXCLUDED.context,
			confidence = EXCLUDED.confidence,
			attention_score = EXCLUDED.attention_score,
			context_type = EXCLUDED.context_type,
			attribution = EXCLUDED.attribution,
			attribution_confidence = EXCLUDED.attribution_confidence,
			attribution_source = EXCLUDED.attribution_source
	`,
		instrumentID,
		c.Type,
		c.Severity,
		c.PreviousValue,
		c.CurrentValue,
		c.ChangePercent,
		c.CurrentVolume,
		c.AverageVolume,
		c.VolumeRatio,
		c.VolumeAnomaly,
		c.PriceAnomaly,
		c.Near52WeekHigh,
		c.Near52WeekLow,
		strings.Join(c.Signals, ","),
		c.Context,
		c.Confidence,
		c.AttentionScore,
		c.ContextType,
		c.Attribution,
		c.AttributionConfidence,
		c.AttributionSource,
		c.DetectedAt,
	)

	return err
}

func (r *Repository) Since(
	ctx context.Context,
	instrumentID string,
	since time.Time,
) ([]MarketChange, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			instrument_id,
			change_type,
			severity,
			previous_value,
			current_value,
			change_percent,
			current_volume,
			average_volume,
			volume_ratio,
			volume_anomaly,
			price_anomaly,
			near_52w_high,
			near_52w_low,
			signals,
			context,
			confidence,
			attention_score,
			detected_at
		FROM market_changes
		WHERE instrument_id = $1
		  AND detected_at > $2
		ORDER BY detected_at ASC
	`,
		instrumentID,
		since,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]MarketChange, 0)

	for rows.Next() {

		var c MarketChange
		var signals string

		if err := rows.Scan(
			&c.ID,
			&c.InstrumentID,
			&c.Type,
			&c.Severity,
			&c.PreviousValue,
			&c.CurrentValue,
			&c.ChangePercent,
			&c.CurrentVolume,
			&c.AverageVolume,
			&c.VolumeRatio,
			&c.VolumeAnomaly,
			&c.PriceAnomaly,
			&c.Near52WeekHigh,
			&c.Near52WeekLow,
			&signals,
			&c.Context,
			&c.Confidence,
			&c.AttentionScore,
			&c.DetectedAt,
		); err != nil {
			return nil, err
		}

		if signals != "" {
			c.Signals = strings.Split(signals, ",")
		}

		result = append(result, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) History(
	ctx context.Context,
	instrumentID string,
	limit int,
) ([]MarketChange, error) {

	if limit <= 0 {
		limit = 20
	}

	rows, err := r.db.QueryContext(ctx, `
        SELECT
            mc.id,
            mc.instrument_id,
            i.symbol,
            mc.change_type,
            mc.severity,
            mc.previous_value,
            mc.current_value,
            mc.change_percent,
            mc.current_volume,
            mc.average_volume,
            mc.volume_ratio,
            mc.volume_anomaly,
            mc.price_anomaly,
            mc.near_52w_high,
            mc.near_52w_low,
            mc.signals,
            mc.context,
            mc.confidence,
            mc.attention_score,
            mc.context_type,
            mc.attribution,
            mc.attribution_confidence,
            mc.attribution_source,
            mc.detected_at
        FROM market_changes mc
        JOIN instruments i
            ON i.id = mc.instrument_id
        WHERE mc.instrument_id = $1
        ORDER BY mc.detected_at DESC
        LIMIT $2
    `, instrumentID, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	changes := make([]MarketChange, 0)

	for rows.Next() {
		var c MarketChange
		var signals string

		err := rows.Scan(
			&c.ID,
			&c.InstrumentID,
			&c.Symbol,
			&c.Type,
			&c.Severity,
			&c.PreviousValue,
			&c.CurrentValue,
			&c.ChangePercent,
			&c.CurrentVolume,
			&c.AverageVolume,
			&c.VolumeRatio,
			&c.VolumeAnomaly,
			&c.PriceAnomaly,
			&c.Near52WeekHigh,
			&c.Near52WeekLow,
			&signals,
			&c.Context,
			&c.Confidence,
			&c.AttentionScore,
			&c.ContextType,
			&c.Attribution,
			&c.AttributionConfidence,
			&c.AttributionSource,
			&c.DetectedAt,
		)

		if err != nil {
			return nil, err
		}

		if signals != "" {
			c.Signals = strings.Split(signals, ",")
		}

		changes = append(changes, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return changes, nil
}
