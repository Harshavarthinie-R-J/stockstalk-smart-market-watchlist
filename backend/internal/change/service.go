package change

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Detect creates the basic price event.
func (s *Service) Detect(
	previous,
	current,
	weekHigh,
	weekLow float64,
) MarketChange {

	percent := 0.0

	if previous != 0 {
		percent = ((current - previous) / previous) * 100
	}

	magnitude := math.Abs(percent)

	severity := Normal

	if magnitude >= 5 {
		severity = Significant
	} else if magnitude >= 2 {
		severity = Notable
	}

	changeType := PriceMovement

	nearHigh := false
	nearLow := false

	if weekHigh > 0 && current >= weekHigh {
		changeType = NewHigh
		nearHigh = true
	} else if weekLow > 0 && current <= weekLow {
		changeType = NewLow
		nearLow = true
	} else {
		// Treat within 2% of the 52-week extreme as "near".
		if weekHigh > 0 && current >= weekHigh*0.98 {
			nearHigh = true
		}

		if weekLow > 0 && current <= weekLow*1.02 {
			nearLow = true
		}
	}

	change := MarketChange{
		ID:             uuid.NewString(),
		Type:           changeType,
		Severity:       severity,
		PreviousValue:  previous,
		CurrentValue:   current,
		ChangePercent:  percent,
		Near52WeekHigh: nearHigh,
		Near52WeekLow:  nearLow,
		PriceAnomaly:   magnitude >= 2,
		DetectedAt:     time.Now().UTC(),
	}

	change.Signals = append(change.Signals,
		"PRICE_MOVEMENT",
	)

	if nearHigh {
		change.Signals = append(change.Signals, "NEAR_52W_HIGH")
	}

	if nearLow {
		change.Signals = append(change.Signals, "NEAR_52W_LOW")
	}

	change.AttentionScore = calculateAttentionScore(change)

	return change
}

// EnrichWithVolume adds personal volume baseline information.
func (s *Service) EnrichWithVolume(
	ctx context.Context,
	c *MarketChange,
	currentVolume int64,
	averageVolume float64,
) {
	c.CurrentVolume = currentVolume
	c.AverageVolume = averageVolume

	if averageVolume <= 0 {
		c.VolumeRatio = 0
		c.VolumeAnomaly = false
		c.AttentionScore = calculateAttentionScore(*c)
		return
	}

	c.VolumeRatio = float64(currentVolume) / averageVolume

	// 1.5x normal volume = unusual.
	if c.VolumeRatio >= 1.5 {
		c.VolumeAnomaly = true

		// Prevent duplicate signal insertion.
		hasSignal := false
		for _, signal := range c.Signals {
			if signal == "UNUSUAL_VOLUME" {
				hasSignal = true
				break
			}
		}

		if !hasSignal {
			c.Signals = append(c.Signals, "UNUSUAL_VOLUME")
		}
	}

	c.AttentionScore = calculateAttentionScore(*c)
}

// Fuse creates one coherent event from multiple signals.
func (s *Service) Fuse(c MarketChange) MarketChange {

	signalCount := len(c.Signals)

	if signalCount >= 2 {
		c.Type = FusedEvent
	}

	// Context is deliberately factual rather than claiming causation.
	switch {
	case c.VolumeAnomaly && c.Near52WeekHigh:
		c.Context = "Price movement occurred with unusual volume near the 52-week high."
		c.Confidence = HighConfidence

	case c.VolumeAnomaly && c.Near52WeekLow:
		c.Context = "Price movement occurred with unusual volume near the 52-week low."
		c.Confidence = HighConfidence

	case c.VolumeAnomaly:
		c.Context = "Price movement occurred with unusually high trading volume."
		c.Confidence = MediumConfidence

	case c.Near52WeekHigh:
		c.Context = "Price is near the 52-week high."
		c.Confidence = MediumConfidence

	case c.Near52WeekLow:
		c.Context = "Price is near the 52-week low."
		c.Confidence = MediumConfidence

	default:
		c.Context = "Price movement detected, but no additional strong context was found."
		c.Confidence = NoClearContext
	}

	c.AttentionScore = calculateAttentionScore(c)

	return c
}

// AddAttribution adds structured context and attribution information.
// Attribution is based only on observed market data.
// It does not claim that one factor caused another.
func (s *Service) AddAttribution(c MarketChange) MarketChange {

	switch {
	case c.VolumeAnomaly && c.Near52WeekHigh:
		c.ContextType = "UNUSUAL_VOLUME_NEAR_52W_HIGH"
		c.Attribution = "Unusual trading volume was observed while the price was near its 52-week high."
		c.AttributionConfidence = HighConfidence
		c.AttributionSource = "MARKET_OBSERVATION"

	case c.VolumeAnomaly && c.Near52WeekLow:
		c.ContextType = "UNUSUAL_VOLUME_NEAR_52W_LOW"
		c.Attribution = "Unusual trading volume was observed while the price was near its 52-week low."
		c.AttributionConfidence = HighConfidence
		c.AttributionSource = "MARKET_OBSERVATION"

	case c.VolumeAnomaly:
		c.ContextType = "UNUSUAL_VOLUME"
		c.Attribution = "Trading volume is unusually high compared with the historical baseline."
		c.AttributionConfidence = MediumConfidence
		c.AttributionSource = "MARKET_OBSERVATION"

	case c.Near52WeekHigh:
		c.ContextType = "NEAR_52W_HIGH"
		c.Attribution = "The current price is close to the stock's 52-week high."
		c.AttributionConfidence = MediumConfidence
		c.AttributionSource = "MARKET_OBSERVATION"

	case c.Near52WeekLow:
		c.ContextType = "NEAR_52W_LOW"
		c.Attribution = "The current price is close to the stock's 52-week low."
		c.AttributionConfidence = MediumConfidence
		c.AttributionSource = "MARKET_OBSERVATION"

	case c.PriceAnomaly:
		c.ContextType = "PRICE_MOVEMENT"
		c.Attribution = "A significant price movement was observed."
		c.AttributionConfidence = MediumConfidence
		c.AttributionSource = "MARKET_OBSERVATION"

	default:
		c.ContextType = "NO_CLEAR_CONTEXT"
		c.Attribution = "No additional strong market context was observed."
		c.AttributionConfidence = NoClearContext
		c.AttributionSource = "MARKET_OBSERVATION"
	}

	return c
}

func (s *Service) History(
	ctx context.Context,
	instrumentID string,
	limit int,
) ([]MarketChange, error) {
	return s.repo.History(ctx, instrumentID, limit)
}

// calculateAttentionScore gives every event an explainable score from 0-100.
func calculateAttentionScore(c MarketChange) int {

	score := 0

	magnitude := math.Abs(c.ChangePercent)

	// Price movement: 0-40
	switch {
	case magnitude >= 10:
		score += 40
	case magnitude >= 5:
		score += 30
	case magnitude >= 2:
		score += 20
	default:
		score += int(magnitude * 10)
	}

	// Severity: 5-25
	switch c.Severity {
	case Significant:
		score += 25
	case Notable:
		score += 15
	case Normal:
		score += 5
	}

	// Volume anomaly: 0-20
	switch {
	case c.VolumeRatio >= 3:
		score += 20
	case c.VolumeRatio >= 2:
		score += 15
	case c.VolumeRatio >= 1.5:
		score += 10
	}

	// 52-week event: 20
	if c.Type == NewHigh || c.Type == NewLow {
		score += 20
	}

	// Near 52-week extreme: 10
	if c.Near52WeekHigh || c.Near52WeekLow {
		score += 10
	}

	if score > 100 {
		score = 100
	}

	if score < 0 {
		score = 0
	}

	return score
}

func (s *Service) Save(
	ctx context.Context,
	instrumentID string,
	c MarketChange,
) error {
	return s.repo.Create(ctx, instrumentID, c)
}
