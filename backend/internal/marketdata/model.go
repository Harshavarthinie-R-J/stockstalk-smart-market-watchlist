package marketdata

import "time"

type Instrument struct {
	ID       string `json:"id"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Segment  string `json:"segment"`
	ISIN     string `json:"isin"`
	Currency string `json:"currency"`
	Sector   string `json:"sector"`
	Industry string `json:"industry"`
}

type Quote struct {
	InstrumentID  string    `json:"instrument_id"`
	Symbol        string    `json:"symbol"`
	Price         float64   `json:"price"`
	PreviousClose float64   `json:"previous_close"`
	Open          float64   `json:"open"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
	Volume        int64     `json:"volume"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"change_percent"`
	Week52High    float64   `json:"week_52_high"`
	Week52Low     float64   `json:"week_52_low"`
	MarketStatus  string    `json:"market_status"`
	MarketTime    time.Time `json:"market_timestamp"`
	ReceivedTime  time.Time `json:"received_timestamp"`
	Source        string    `json:"source"`
}

type MarketData struct {
	Market      string       `json:"market"`
	LastUpdated time.Time    `json:"last_updated"`
	Instruments []Instrument `json:"instruments"`
	Quotes      []Quote      `json:"quotes"`
}
