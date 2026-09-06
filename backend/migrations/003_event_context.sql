ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS context_type TEXT NOT NULL DEFAULT 'NO_CLEAR_CONTEXT';

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS attribution TEXT NOT NULL DEFAULT '';

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS attribution_confidence TEXT NOT NULL DEFAULT 'NO_CLEAR_CONTEXT';

ALTER TABLE market_changes
ADD COLUMN IF NOT EXISTS attribution_source TEXT NOT NULL DEFAULT 'MARKET_OBSERVATION';

CREATE INDEX IF NOT EXISTS idx_market_changes_context
ON market_changes(context_type);

CREATE INDEX IF NOT EXISTS idx_market_changes_attention
ON market_changes(attention_score DESC);