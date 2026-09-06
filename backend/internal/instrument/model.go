package instrument

type Instrument struct {
	ID       string
	Symbol   string
	Name     string
	Exchange string
	Segment  string
	ISIN     string
	Currency string
	Sector   string
	Industry string
	Active   bool
}
