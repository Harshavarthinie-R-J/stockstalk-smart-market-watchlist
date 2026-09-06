CREATE UNIQUE INDEX IF NOT EXISTS idx_market_snapshots_dedup
ON market_snapshots(instrument_id, market_timestamp, source);

CREATE INDEX IF NOT EXISTS idx_market_snapshots_history
ON market_snapshots(instrument_id, market_timestamp ASC);