package change

import "time"

type Type string

const (
	PriceMovement Type = "PRICE_MOVEMENT"
	NewHigh       Type = "NEW_52W_HIGH"
	NewLow        Type = "NEW_52W_LOW"
	FusedEvent    Type = "FUSED_MARKET_EVENT"
)

type Severity string

const (
	Normal      Severity = "NORMAL"
	Notable     Severity = "NOTABLE"
	Significant Severity = "SIGNIFICANT"
)

type Confidence string

const (
	HighConfidence   Confidence = "HIGH"
	MediumConfidence Confidence = "MEDIUM"
	LowConfidence    Confidence = "LOW"
	NoClearContext   Confidence = "NO_CLEAR_CONTEXT"
)

type MarketChange struct {
	ID           string
	InstrumentID string
	Symbol       string
	Type         Type
	Severity     Severity

	PreviousValue float64
	CurrentValue  float64
	ChangePercent float64

	// Market intelligence
	CurrentVolume  int64
	AverageVolume  float64
	VolumeRatio    float64
	VolumeAnomaly  bool
	PriceAnomaly   bool
	Near52WeekHigh bool
	Near52WeekLow  bool

	// Event explanation
	Signals               []string
	Context               string
	Confidence            Confidence
	ContextType           string
	Attribution           string
	AttributionConfidence Confidence
	AttributionSource     string
	// Ranking
	AttentionScore int

	DetectedAt       time.Time
	ReliabilityState string
}
