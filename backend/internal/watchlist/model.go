package watchlist

import "time"

type Watchlist struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Stocks    []WatchItem `json:"stocks"`
}

type WatchItem struct {
	InstrumentID string    `json:"instrument_id"`
	Position     int       `json:"position"`
	AddedAt      time.Time `json:"added_at"`
}
