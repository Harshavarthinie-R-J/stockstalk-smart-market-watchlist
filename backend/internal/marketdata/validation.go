package marketdata

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

func ValidateQuote(q Quote) error {
	if strings.TrimSpace(q.InstrumentID) == "" {
		return errors.New("instrument id is required")
	}

	if strings.TrimSpace(q.Symbol) == "" {
		return errors.New("symbol is required")
	}

	if q.Price <= 0 || math.IsNaN(q.Price) || math.IsInf(q.Price, 0) {
		return errors.New("price must be positive and finite")
	}

	if q.PreviousClose <= 0 {
		return errors.New("previous close must be positive")
	}

	if q.Open <= 0 {
		return errors.New("open price must be positive")
	}

	if q.High <= 0 || q.Low <= 0 {
		return errors.New("high and low must be positive")
	}

	if q.Low > q.High {
		return errors.New("low price cannot exceed high price")
	}

	if q.Price < q.Low || q.Price > q.High {
		return fmt.Errorf(
			"price %.4f is outside high/low range %.4f-%.4f",
			q.Price,
			q.Low,
			q.High,
		)
	}

	if q.Volume < 0 {
		return errors.New("volume cannot be negative")
	}

	if q.Week52Low < 0 || q.Week52High < 0 {
		return errors.New("52-week values cannot be negative")
	}

	if q.Week52High > 0 &&
		q.Week52Low > 0 &&
		q.Week52Low > q.Week52High {
		return errors.New("52-week low cannot exceed 52-week high")
	}

	if q.MarketTime.IsZero() {
		return errors.New("market timestamp is required")
	}

	if q.MarketTime.After(time.Now().Add(5 * time.Minute)) {
		return errors.New("future market timestamp")
	}

	if strings.TrimSpace(q.Source) == "" {
		return errors.New("data source is required")
	}

	return nil
}
