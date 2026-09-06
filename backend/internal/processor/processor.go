package processor

import (
	"context"
	"stockstalk/internal/change"
	"stockstalk/internal/marketdata"
	"stockstalk/internal/watchlist"
	"time"
)

type Processor struct {
	Market           *marketdata.Service
	SnapshotRepo     *marketdata.Repository
	WatchlistService *watchlist.Service
	ChangeService    *change.Service
	Publish          func(userID, watchlistID string, event interface{})
}

func New(
	market *marketdata.Service,
	snapshotRepo *marketdata.Repository,
	watchlistService *watchlist.Service,
	changeService *change.Service,
	publish func(userID, watchlistID string, event interface{}),
) *Processor {
	return &Processor{
		Market:           market,
		SnapshotRepo:     snapshotRepo,
		WatchlistService: watchlistService,
		ChangeService:    changeService,
		Publish:          publish,
	}
}

// Process fetches the newest market state, stores new snapshots,
// detects meaningful changes, enriches them, and publishes events.
func (p *Processor) Process(ctx context.Context) error {
	data, err := p.Market.GetMarketData()
	if err != nil {
		return err
	}

	watchlists, err := p.WatchlistService.ListAll(ctx)
	if err != nil {
		return err
	}

	// Save each new market snapshot first.
	for _, quote := range data.Quotes {
		latest, err := p.SnapshotRepo.LatestSnapshot(
			ctx,
			quote.InstrumentID,
		)

		if err == nil &&
			latest.MarketTime.Equal(quote.MarketTime) &&
			latest.Source == quote.Source {
			// Nothing new for this instrument.
			continue
		}

		if err := p.SnapshotRepo.SaveSnapshot(ctx, quote); err != nil {
			continue
		}

		previous, err := p.SnapshotRepo.SnapshotBefore(
			ctx,
			quote.InstrumentID,
			quote.MarketTime,
		)

		if err != nil || previous == nil {
			continue
		}

		detected := p.ChangeService.Detect(
			previous.Price,
			quote.Price,
			quote.Week52High,
			quote.Week52Low,
		)

		detected.InstrumentID = quote.InstrumentID
		detected.Symbol = quote.Symbol
		detected.ReliabilityState = marketdata.ReliabilityState(&quote)

		averageVolume, err := p.SnapshotRepo.AverageVolume(
			ctx,
			quote.InstrumentID,
			quote.MarketTime,
		)

		if err == nil {
			p.ChangeService.EnrichWithVolume(
				ctx,
				&detected,
				quote.Volume,
				averageVolume,
			)
		}

		detected = p.ChangeService.Fuse(detected)
		detected = p.ChangeService.AddAttribution(detected)

		if detected.Severity == change.Normal {
			continue
		}

		// Only send the event to watchlists containing this instrument.
		for _, w := range watchlists {
			contains := false

			for _, stock := range w.Stocks {
				if stock.InstrumentID == quote.InstrumentID {
					contains = true
					break
				}
			}

			if !contains {
				continue
			}

			// Store an event for this watchlist processing.
			if err := p.ChangeService.Save(ctx, w.ID, detected); err != nil {
				continue
			}

			if p.Publish != nil {
				p.Publish(
					w.UserID,
					w.ID,
					changeEventMap(&detected),
				)
			}
		}
	}

	return nil
}

func changeEventMap(c *change.MarketChange) map[string]interface{} {
	return map[string]interface{}{
		"id":                    c.ID,
		"instrumentId":          c.InstrumentID,
		"symbol":                c.Symbol,
		"type":                  string(c.Type),
		"severity":              string(c.Severity),
		"previousValue":         c.PreviousValue,
		"currentValue":          c.CurrentValue,
		"changePercent":         c.ChangePercent,
		"currentVolume":         c.CurrentVolume,
		"averageVolume":         c.AverageVolume,
		"volumeRatio":           c.VolumeRatio,
		"volumeAnomaly":         c.VolumeAnomaly,
		"priceAnomaly":          c.PriceAnomaly,
		"near52WeekHigh":        c.Near52WeekHigh,
		"near52WeekLow":         c.Near52WeekLow,
		"signals":               c.Signals,
		"context":               c.Context,
		"confidence":            string(c.Confidence),
		"contextType":           c.ContextType,
		"attribution":           c.Attribution,
		"attributionConfidence": string(c.AttributionConfidence),
		"attributionSource":     c.AttributionSource,
		"attentionScore":        c.AttentionScore,
		"reliabilityState":      c.ReliabilityState,
		"detectedAt":            c.DetectedAt.Format(time.RFC3339),
	}
}
