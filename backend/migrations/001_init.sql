CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS instruments (
    id TEXT PRIMARY KEY,
    symbol TEXT NOT NULL,
    name TEXT NOT NULL,
    exchange TEXT NOT NULL,
    segment TEXT NOT NULL,
    isin TEXT,
    currency TEXT NOT NULL,
    sector TEXT,
    industry TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_instruments_symbol ON instruments(symbol);

CREATE TABLE IF NOT EXISTS watchlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name)
);

CREATE TABLE IF NOT EXISTS watchlist_stocks (
    watchlist_id UUID NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
    instrument_id TEXT NOT NULL REFERENCES instruments(id) ON DELETE RESTRICT,
    position INTEGER NOT NULL DEFAULT 0,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (watchlist_id, instrument_id)
);

CREATE INDEX IF NOT EXISTS idx_watchlist_stocks_watchlist ON watchlist_stocks(watchlist_id, position);

CREATE TABLE IF NOT EXISTS market_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instrument_id TEXT NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    price NUMERIC(20,8) NOT NULL,
    previous_close NUMERIC(20,8) NOT NULL,
    open_price NUMERIC(20,8) NOT NULL,
    high_price NUMERIC(20,8) NOT NULL,
    low_price NUMERIC(20,8) NOT NULL,
    volume BIGINT NOT NULL,
    change_value NUMERIC(20,8) NOT NULL,
    change_percent NUMERIC(12,6) NOT NULL,
    week_52_high NUMERIC(20,8) NOT NULL,
    week_52_low NUMERIC(20,8) NOT NULL,
    market_status TEXT NOT NULL,
    market_timestamp TIMESTAMPTZ NOT NULL,
    received_timestamp TIMESTAMPTZ NOT NULL,
    source TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_market_snapshots_instrument_time
    ON market_snapshots(instrument_id, market_timestamp DESC);

CREATE TABLE IF NOT EXISTS checkpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    watchlist_id UUID NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_checkpoints_lookup
    ON checkpoints(user_id, watchlist_id, created_at DESC);

CREATE TABLE IF NOT EXISTS market_changes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instrument_id TEXT NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    change_type TEXT NOT NULL,
    severity TEXT NOT NULL,
    previous_value NUMERIC(20,8) NOT NULL,
    current_value NUMERIC(20,8) NOT NULL,
    change_percent NUMERIC(12,6) NOT NULL,
    detected_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_market_changes_instrument_time
    ON market_changes(instrument_id, detected_at DESC);
