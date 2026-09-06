package marketdata

import (
	"testing"
	"time"
)

func TestReliabilityState(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		quote     *Quote
		wantState string
	}{
		{
			name:      "unavailable nil quote",
			quote:     nil,
			wantState: "UNAVAILABLE",
		},
		{
			name: "unavailable zero timestamp",
			quote: &Quote{
				MarketTime: time.Time{},
			},
			wantState: "UNAVAILABLE",
		},
		{
			name: "live",
			quote: &Quote{
				MarketTime: now.Add(-1 * time.Minute),
			},
			wantState: "LIVE",
		},
		{
			name: "delayed",
			quote: &Quote{
				MarketTime: now.Add(-10 * time.Minute),
			},
			wantState: "DELAYED",
		},
		{
			name: "stale",
			quote: &Quote{
				MarketTime: now.Add(-2 * time.Hour),
			},
			wantState: "STALE",
		},
		{
			name: "conflicting future timestamp",
			quote: &Quote{
				MarketTime: now.Add(10 * time.Minute),
			},
			wantState: "CONFLICTING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReliabilityState(tt.quote)

			if got != tt.wantState {
				t.Fatalf(
					"state = %s, want %s",
					got,
					tt.wantState,
				)
			}
		})
	}
}
