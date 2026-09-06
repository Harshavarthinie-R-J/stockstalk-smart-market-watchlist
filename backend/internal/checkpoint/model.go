package checkpoint

import "time"

type Checkpoint struct {
	ID          string
	UserID      string
	WatchlistID string
	CreatedAt   time.Time
}
