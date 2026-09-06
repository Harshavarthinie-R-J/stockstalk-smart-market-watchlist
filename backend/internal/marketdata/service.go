package marketdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Service struct {
	data MarketData
}

func NewService(filePath string) (*Service, error) {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var raw struct {
		Market      string       `json:"market"`
		LastUpdated string       `json:"last_updated"`
		Instruments []Instrument `json:"instruments"`
		Quotes      []struct {
			InstrumentID  string  `json:"instrument_id"`
			Symbol        string  `json:"symbol"`
			Price         float64 `json:"price"`
			PreviousClose float64 `json:"previous_close"`
			Open          float64 `json:"open"`
			High          float64 `json:"high"`
			Low           float64 `json:"low"`
			Volume        int64   `json:"volume"`
			Change        float64 `json:"change"`
			ChangePercent float64 `json:"change_percent"`
			Week52High    float64 `json:"week_52_high"`
			Week52Low     float64 `json:"week_52_low"`
			MarketStatus  string  `json:"market_status"`
			MarketTime    string  `json:"market_timestamp"`
			Source        string  `json:"source"`
		} `json:"quotes"`
	}

	if err := json.Unmarshal(contents, &raw); err != nil {
		return nil, err
	}

	lastUpdated, err := time.Parse(time.RFC3339, raw.LastUpdated)
	if err != nil {
		return nil, fmt.Errorf("invalid last_updated timestamp: %w", err)
	}

	result := MarketData{
		Market:      raw.Market,
		LastUpdated: lastUpdated,
		Instruments: raw.Instruments,
	}

	now := time.Now().UTC()

	for _, q := range raw.Quotes {
		marketTime, err := time.Parse(time.RFC3339, q.MarketTime)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid market timestamp for %s: %w",
				q.Symbol,
				err,
			)
		}

		quote := Quote{
			InstrumentID:  q.InstrumentID,
			Symbol:        q.Symbol,
			Price:         q.Price,
			PreviousClose: q.PreviousClose,
			Open:          q.Open,
			High:          q.High,
			Low:           q.Low,
			Volume:        q.Volume,
			Change:        q.Change,
			ChangePercent: q.ChangePercent,
			Week52High:    q.Week52High,
			Week52Low:     q.Week52Low,
			MarketStatus:  q.MarketStatus,
			MarketTime:    marketTime,
			ReceivedTime:  now,
			Source:        q.Source,
		}

		// Reject malformed or impossible market observations.
		if err := ValidateQuote(quote); err != nil {
			return nil, fmt.Errorf(
				"invalid quote for %s: %w",
				q.Symbol,
				err,
			)
		}

		result.Quotes = append(result.Quotes, quote)
	}

	return &Service{data: result}, nil
}

func (s *Service) GetMarketData() (*MarketData, error) {
	return &s.data, nil
}

func (s *Service) GetQuote(symbol string) (*Quote, error) {
	for _, quote := range s.data.Quotes {
		if strings.EqualFold(quote.Symbol, symbol) {
			result := quote
			return &result, nil
		}
	}

	return nil, errors.New("quote not found")
}

func (s *Service) GetInstrument(symbol string) (*Instrument, error) {
	for _, item := range s.data.Instruments {
		if strings.EqualFold(item.Symbol, symbol) {
			result := item
			return &result, nil
		}
	}

	return nil, errors.New("instrument not found")
}

// ReliabilityState describes how trustworthy the current market data is.
func ReliabilityState(q *Quote) string {
	if q == nil {
		return "UNAVAILABLE"
	}

	if q.MarketTime.IsZero() {
		return "UNAVAILABLE"
	}

	age := time.Since(q.MarketTime)

	// Future timestamps are treated as conflicting data.
	if age < 0 {
		return "CONFLICTING"
	}

	switch {
	case age <= 5*time.Minute:
		return "LIVE"

	case age <= 30*time.Minute:
		return "DELAYED"

	default:
		return "STALE"
	}
}
