ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS current_volume BIGINT NOT NULL DEFAULT 0;

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS average_volume NUMERIC(20,8) NOT NULL DEFAULT 0;

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS volume_ratio NUMERIC(20,8) NOT NULL DEFAULT 1;

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS volume_anomaly BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS price_anomaly BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS near_52w_high BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS near_52w_low BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS signals TEXT NOT NULL DEFAULT '';

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS context TEXT NOT NULL DEFAULT '';

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS confidence TEXT NOT NULL DEFAULT 'NO_CLEAR_CONTEXT';

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS attention_score INTEGER NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX IF NOT EXISTS idx_market_changes_dedup
ON market_changes(
    instrument_id,
    change_type,
    previous_value,
    current_value,
    change_percent
);