package marketdata

type Provider interface {
	GetMarketData() (*MarketData, error)
	GetQuote(symbol string) (*Quote, error)
	GetInstrument(symbol string) (*Instrument, error)
}
