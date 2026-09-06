package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"stockstalk/internal/auth"
	"stockstalk/internal/change"
	"stockstalk/internal/checkpoint"
	"stockstalk/internal/instrument"
	"stockstalk/internal/marketdata"
	"stockstalk/internal/user"
	"stockstalk/internal/watchlist"
	"time"
	"net/http"
)

type Resolver struct {
	UserService       *user.Service
	AuthService       *auth.Service
	WatchlistService  *watchlist.Service
	InstrumentService *instrument.Service
	MarketService     *marketdata.Service
	SnapshotRepo      *marketdata.Repository
	CheckpointService *checkpoint.Service
	ChangeService     *change.Service
	Broker            *eventBroker
}

func NewResolver(
	userService *user.Service,
	authService *auth.Service,
	watchlistService *watchlist.Service,
	instrumentService *instrument.Service,
	marketService *marketdata.Service,
	snapshotRepo *marketdata.Repository,
	checkpointService *checkpoint.Service,
	changeService *change.Service,

) *Resolver {
	return &Resolver{
		UserService:       userService,
		AuthService:       authService,
		WatchlistService:  watchlistService,
		InstrumentService: instrumentService,
		MarketService:     marketService,
		SnapshotRepo:      snapshotRepo,
		CheckpointService: checkpointService,
		ChangeService:     changeService,
		Broker:            newEventBroker(),
	}
}

// Auth guard.
func (r *Resolver) userID(ctx context.Context) (string, error) {
	id, ok := auth.UserID(ctx)
	if !ok {
		return "", errors.New("authentication required")
	}

	return id, nil
}

// Register.
func (r *Resolver) Register(
	ctx context.Context,
	args map[string]interface{},
) (map[string]interface{}, error) {

	name, ok := args["name"].(string)
	if !ok {
		return nil, errors.New("invalid name")
	}

	email, ok := args["email"].(string)
	if !ok {
		return nil, errors.New("invalid email")
	}

	password, ok := args["password"].(string)
	if !ok {
		return nil, errors.New("invalid password")
	}

	u, token, err := r.AuthService.Register(
		ctx,
		name,
		email,
		password,
	)
	if err != nil {
		return nil, err
	}

	return userMap(u, token), nil
}

// Login.
func (r *Resolver) Login(
	ctx context.Context,
	args map[string]interface{},
) (map[string]interface{}, error) {

	email, ok := args["email"].(string)
	if !ok {
		return nil, errors.New("invalid email")
	}

	password, ok := args["password"].(string)
	if !ok {
		return nil, errors.New("invalid password")
	}

	u, token, err := r.AuthService.Login(
		ctx,
		email,
		password,
	)
	if err != nil {
		return nil, err
	}

	return userMap(u, token), nil
}

// Current user.
func (r *Resolver) Me(
	ctx context.Context,
) (map[string]interface{}, error) {

	id, err := r.userID(ctx)
	if err != nil {
		return nil, err
	}

	u, err := r.UserService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return userMap(u, ""), nil
}

// Market data.
func (r *Resolver) Market(
	ctx context.Context,
) (map[string]interface{}, error) {

	data, err := r.MarketService.GetMarketData()
	if err != nil {
		return nil, err
	}

	return marketMap(data), nil
}

// Single quote.
func (r *Resolver) Quote(
	ctx context.Context,
	args map[string]interface{},
) (map[string]interface{}, error) {

	symbol, ok := args["symbol"].(string)
	if !ok {
		return nil, errors.New("invalid symbol")
	}

	q, err := r.MarketService.GetQuote(symbol)
	if err != nil {
		return nil, err
	}

	return quoteMap(q), nil
}

// Instrument search.
func (r *Resolver) SearchInstruments(
	ctx context.Context,
	args map[string]interface{},
) ([]map[string]interface{}, error) {

	query, ok := args["query"].(string)
	if !ok {
		return nil, errors.New("invalid search query")
	}

	items, err := r.InstrumentService.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(items))

	for _, item := range items {
		result = append(
			result,
			instrumentMap(&item),
		)
	}

	return result, nil
}

// User watchlists.
func (r *Resolver) Watchlists(
	ctx context.Context,
) ([]map[string]interface{}, error) {

	id, err := r.userID(ctx)
	if err != nil {
		return nil, err
	}

	items, err := r.WatchlistService.List(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(items))

	for _, item := range items {
		result = append(
			result,
			watchlistMap(&item),
		)
	}

	return result, nil
}

// Single watchlist.
func (r *Resolver) Watchlist(
	ctx context.Context,
	args map[string]interface{},
) (map[string]interface{}, error) {

	id, err := r.userID(ctx)
	if err != nil {
		return nil, err
	}

	watchlistID, ok := args["id"].(string)
	if !ok {
		return nil, errors.New("invalid watchlist id")
	}

	w, err := r.WatchlistService.Get(
		ctx,
		id,
		watchlistID,
	)
	if err != nil {
		return nil, err
	}

	return watchlistMap(w), nil
}

// Create watchlist.
func (r *Resolver) CreateWatchlist(
	ctx context.Context,
	args map[string]interface{},
) (map[string]interface{}, error) {

	id, err := r.userID(ctx)
	if err != nil {
		return nil, err
	}

	name, ok := args["name"].(string)
	if !ok {
		return nil, errors.New("invalid watchlist name")
	}

	w, err := r.WatchlistService.Create(
		ctx,
		id,
		name,
	)
	if err != nil {
		return nil, err
	}

	return watchlistMap(w), nil
}

// Add stock.
func (r *Resolver) AddStock(
	ctx context.Context,
	args map[string]interface{},
) (bool, error) {

	id, err := r.userID(ctx)
	if err != nil {
		return false, err
	}

	watchlistID, ok := args["watchlistId"].(string)
	if !ok {
		return false, errors.New("invalid watchlist id")
	}

	instrumentID, ok := args["instrumentId"].(string)
	if !ok {
		return false, errors.New("invalid instrument id")
	}

	if _, err := r.InstrumentService.GetByID(
		ctx,
		instrumentID,
	); err != nil {
		return false, errors.New("instrument not found")
	}

	err = r.WatchlistService.AddStock(
		ctx,
		id,
		watchlistID,
		instrumentID,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Remove stock.
func (r *Resolver) RemoveStock(
	ctx context.Context,
	args map[string]interface{},
) (bool, error) {

	id, err := r.userID(ctx)
	if err != nil {
		return false, err
	}

	watchlistID, ok := args["watchlistId"].(string)
	if !ok {
		return false, errors.New("invalid watchlist id")
	}

	instrumentID, ok := args["instrumentId"].(string)
	if !ok {
		return false, errors.New("invalid instrument id")
	}

	err = r.WatchlistService.RemoveStock(
		ctx,
		id,
		watchlistID,
		instrumentID,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Create checkpoint.
func (r *Resolver) CreateCheckpoint(
	ctx context.Context,
	args map[string]interface{},
) (map[string]interface{}, error) {

	id, err := r.userID(ctx)
	if err != nil {
		return nil, err
	}

	watchlistID, ok := args["watchlistId"].(string)
	if !ok {
		return nil, errors.New("invalid watchlist id")
	}

	c, err := r.CheckpointService.Create(
		ctx,
		id,
		watchlistID,
	)
	if err != nil {
		return nil, err
	}

	return checkpointMap(c), nil
}

// Changes since checkpoint.
func (r *Resolver) ChangesSinceLastVisit(
	ctx context.Context,
	args map[string]interface{},
) (map[string]interface{}, error) {

	// --------------------------------------------------------
	// 1. Authentication
	// --------------------------------------------------------

	id, err := r.userID(ctx)
	if err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 2. Get watchlist ID
	// --------------------------------------------------------

	watchlistID, ok := args["watchlistId"].(string)
	if !ok || watchlistID == "" {
		return nil, errors.New("invalid watchlist id")
	}

	// --------------------------------------------------------
	// 3. Get watchlist
	// --------------------------------------------------------

	w, err := r.WatchlistService.Get(
		ctx,
		id,
		watchlistID,
	)
	if err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 4. Get latest checkpoint
	// --------------------------------------------------------

	checkpoint, err := r.CheckpointService.Latest(
		ctx,
		id,
		watchlistID,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// --------------------------------------------------------
	// 5. Determine "since" time
	// --------------------------------------------------------

	var since time.Time

	if checkpoint != nil {
		since = checkpoint.CreatedAt
	} else {
		since = time.Time{}
	}

	// --------------------------------------------------------
	// 6. Prepare response
	// --------------------------------------------------------

	changes := make([]map[string]interface{}, 0)

	meaningful := 0
	notable := 0
	unchanged := 0

	// --------------------------------------------------------
	// 7. Process every stock
	// --------------------------------------------------------

	for _, item := range w.Stocks {

		// ----------------------------------------------------
		// 7.1 Get instrument
		// ----------------------------------------------------

		inst, err := r.InstrumentService.GetByID(
			ctx,
			item.InstrumentID,
		)

		if err != nil {
			continue
		}

		// ----------------------------------------------------
		// 7.2 Get current market quote
		// ----------------------------------------------------

		current, err := r.MarketService.GetQuote(
			inst.Symbol,
		)

		if err != nil {
			continue
		}

		// ----------------------------------------------------
		// 7.3 Get previous snapshot
		// ----------------------------------------------------

		previous, err := r.SnapshotRepo.SnapshotAtOrBefore(
			ctx,
			item.InstrumentID,
			since,
		)

		if err != nil {
			continue
		}

		if previous == nil {
			unchanged++
			continue
		}

		// ----------------------------------------------------
		// 7.4 Detect price movement
		// ----------------------------------------------------

		detected := r.ChangeService.Detect(
			previous.Price,
			current.Price,
			current.Week52High,
			current.Week52Low,
		)
		detected.ReliabilityState = marketdata.ReliabilityState(current)

		// ----------------------------------------------------
		// 7.5 Add instrument information
		// ----------------------------------------------------

		detected.InstrumentID = item.InstrumentID
		detected.Symbol = inst.Symbol

		// ----------------------------------------------------
		// 7.6 Calculate historical volume baseline
		// ----------------------------------------------------

		averageVolume, err := r.SnapshotRepo.AverageVolume(
			ctx,
			item.InstrumentID,
			current.MarketTime,
		)
		if err != nil {
			return nil, fmt.Errorf("average volume query failed: %w", err)
		}

		if err != nil {
			return nil, err
		}

		// ----------------------------------------------------
		// 7.7 Detect unusual volume
		// ----------------------------------------------------

		r.ChangeService.EnrichWithVolume(
			ctx,
			&detected,
			current.Volume,
			averageVolume,
		)

		// ----------------------------------------------------
		// 7.8 Fuse detected signals
		// ----------------------------------------------------

		detected = r.ChangeService.Fuse(detected)

		// ----------------------------------------------------
		// 7.9 Add context and attribution
		// ----------------------------------------------------

		detected = r.ChangeService.AddAttribution(detected)

		// ----------------------------------------------------
		// 7.10 Decide whether this stock is interesting
		// ----------------------------------------------------

		if detected.Severity == change.Normal &&
			!detected.VolumeAnomaly &&
			!detected.Near52WeekHigh &&
			!detected.Near52WeekLow {

			unchanged++
			continue
		}

		// ----------------------------------------------------
		// 7.11 Count event importance
		// ----------------------------------------------------

		if detected.Severity == change.Significant {
			meaningful++
		} else {
			notable++
		}

		// ----------------------------------------------------
		// 7.12 Save event to PostgreSQL
		// ----------------------------------------------------

		if err := r.ChangeService.Save(
			ctx,
			item.InstrumentID,
			detected,
		); err != nil {
			return nil, fmt.Errorf("save market change failed: %w", err)
		}
		r.Broker.Publish(
			id,
			watchlistID,
			changeMap(&detected),
		)

		// ----------------------------------------------------
		// 7.13 Debug output
		// ----------------------------------------------------

		fmt.Printf(
			"DEBUG EVENT: symbol=%s since=%s previousPrice=%.2f currentPrice=%.2f currentVolume=%d averageVolume=%.2f volumeRatio=%.2f volumeAnomaly=%v priceAnomaly=%v signals=%v type=%s score=%d contextType=%s attributionConfidence=%s\n",
			inst.Symbol,
			since.Format(time.RFC3339Nano),
			previous.Price,
			current.Price,
			current.Volume,
			averageVolume,
			detected.VolumeRatio,
			detected.VolumeAnomaly,
			detected.PriceAnomaly,
			detected.Signals,
			detected.Type,
			detected.AttentionScore,
			detected.ContextType,
			detected.AttributionConfidence,
		)

		// ----------------------------------------------------
		// 7.14 Add event to GraphQL response
		// ----------------------------------------------------

		changes = append(
			changes,
			changeMap(&detected),
		)
	}

	// --------------------------------------------------------
	// 8. Previous checkpoint information
	// --------------------------------------------------------

	previousCheckpoint := ""

	if checkpoint != nil {
		previousCheckpoint = checkpoint.CreatedAt.Format(
			time.RFC3339,
		)
	}

	// --------------------------------------------------------
	// 9. Return complete summary
	// --------------------------------------------------------

	return map[string]interface{}{
		"watchlistId":        watchlistID,
		"previousCheckpoint": previousCheckpoint,
		"currentTime":        time.Now().UTC().Format(time.RFC3339),
		"totalStocks":        len(w.Stocks),
		"meaningfulChanges":  meaningful,
		"notableChanges":     notable,
		"unchanged":          unchanged,
		"changes":            changes,
	}, nil
}

// User response.
func userMap(
	u *user.User,
	token string,
) map[string]interface{} {

	return map[string]interface{}{
		"id":        u.ID,
		"name":      u.Name,
		"email":     u.Email,
		"createdAt": u.CreatedAt.Format(time.RFC3339),
		"token":     token,
	}
}

// Instrument response.
func instrumentMap(
	i *instrument.Instrument,
) map[string]interface{} {

	return map[string]interface{}{
		"id":       i.ID,
		"symbol":   i.Symbol,
		"name":     i.Name,
		"exchange": i.Exchange,
		"segment":  i.Segment,
		"isin":     i.ISIN,
		"currency": i.Currency,
		"sector":   i.Sector,
		"industry": i.Industry,
	}
}

// Quote response.
func quoteMap(
	q *marketdata.Quote,
) map[string]interface{} {

	return map[string]interface{}{
		"instrumentId":      q.InstrumentID,
		"symbol":            q.Symbol,
		"price":             q.Price,
		"previousClose":     q.PreviousClose,
		"open":              q.Open,
		"high":              q.High,
		"low":               q.Low,
		"volume":            q.Volume,
		"change":            q.Change,
		"changePercent":     q.ChangePercent,
		"week52High":        q.Week52High,
		"week52Low":         q.Week52Low,
		"marketStatus":      q.MarketStatus,
		"marketTimestamp":   q.MarketTime.Format(time.RFC3339),
		"receivedTimestamp": q.ReceivedTime.Format(time.RFC3339),
		"source":            q.Source,
	}
}

// Market response.
func marketMap(
	m *marketdata.MarketData,
) map[string]interface{} {

	instruments := make([]map[string]interface{}, 0, len(m.Instruments))

	for _, item := range m.Instruments {
		instruments = append(
			instruments,
			instrumentMap(&instrument.Instrument{
				ID:       item.ID,
				Symbol:   item.Symbol,
				Name:     item.Name,
				Exchange: item.Exchange,
				Segment:  item.Segment,
				ISIN:     item.ISIN,
				Currency: item.Currency,
				Sector:   item.Sector,
				Industry: item.Industry,
				Active:   true,
			}),
		)
	}

	quotes := make([]map[string]interface{}, 0, len(m.Quotes))

	for i := range m.Quotes {
		quotes = append(
			quotes,
			quoteMap(&m.Quotes[i]),
		)
	}

	return map[string]interface{}{
		"market":      m.Market,
		"lastUpdated": m.LastUpdated.Format(time.RFC3339),
		"instruments": instruments,
		"quotes":      quotes,
	}
}

// Watchlist response.
func watchlistMap(
	w *watchlist.Watchlist,
) map[string]interface{} {

	stocks := make([]map[string]interface{}, 0, len(w.Stocks))

	for _, item := range w.Stocks {
		stocks = append(
			stocks,
			map[string]interface{}{
				"instrumentId": item.InstrumentID,
				"position":     item.Position,
				"addedAt":      item.AddedAt.Format(time.RFC3339),
			},
		)
	}

	return map[string]interface{}{
		"id":        w.ID,
		"name":      w.Name,
		"createdAt": w.CreatedAt.Format(time.RFC3339),
		"updatedAt": w.UpdatedAt.Format(time.RFC3339),
		"stocks":    stocks,
	}
}

// Checkpoint response.
func checkpointMap(
	c *checkpoint.Checkpoint,
) map[string]interface{} {

	return map[string]interface{}{
		"id":          c.ID,
		"watchlistId": c.WatchlistID,
		"createdAt":   c.CreatedAt.Format(time.RFC3339),
	}
}

// Change response.
func changeMap(
	c *change.MarketChange,
) map[string]interface{} {

	return map[string]interface{}{
		"id":            c.ID,
		"instrumentId":  c.InstrumentID,
		"symbol":        c.Symbol,
		"type":          string(c.Type),
		"severity":      string(c.Severity),
		"previousValue": c.PreviousValue,
		"currentValue":  c.CurrentValue,
		"changePercent": c.ChangePercent,

		"currentVolume":  c.CurrentVolume,
		"averageVolume":  c.AverageVolume,
		"volumeRatio":    c.VolumeRatio,
		"volumeAnomaly":  c.VolumeAnomaly,
		"priceAnomaly":   c.PriceAnomaly,
		"near52WeekHigh": c.Near52WeekHigh,
		"near52WeekLow":  c.Near52WeekLow,

		"signals":               c.Signals,
		"context":               c.Context,
		"confidence":            string(c.Confidence),
		"contextType":           c.ContextType,
		"attribution":           c.Attribution,
		"attributionConfidence": string(c.AttributionConfidence),
		"attributionSource":     c.AttributionSource,
		"reliabilityState":      c.ReliabilityState,
		"attentionScore":        c.AttentionScore,
		"detectedAt":            c.DetectedAt.Format(time.RFC3339),
	}
}

// Event history.
func (r *Resolver) EventHistory(
	ctx context.Context,
	args map[string]interface{},
) ([]map[string]interface{}, error) {

	// --------------------------------------------------------
	// 1. Authentication
	// --------------------------------------------------------

	_, err := r.userID(ctx)
	if err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 2. Get instrument ID
	// --------------------------------------------------------

	instrumentID, ok := args["instrumentId"].(string)

	if !ok || instrumentID == "" {
		return nil, errors.New("invalid instrument id")
	}

	// --------------------------------------------------------
	// 3. Get limit
	// --------------------------------------------------------

	limit := 20

	if value, ok := args["limit"]; ok {
		if n, ok := value.(int); ok && n > 0 {
			limit = n
		}
	}

	// --------------------------------------------------------
	// 4. Get history
	// --------------------------------------------------------

	events, err := r.ChangeService.History(
		ctx,
		instrumentID,
		limit,
	)

	if err != nil {
		return nil, err
	}

	// --------------------------------------------------------
	// 5. Convert to GraphQL response
	// --------------------------------------------------------

	result := make([]map[string]interface{}, 0, len(events))

	for i := range events {
		result = append(
			result,
			changeMap(&events[i]),
		)
	}

	return result, nil
}
