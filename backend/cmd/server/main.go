package main

import (
	"context"
	"expvar"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"stockstalk/internal/auth"
	"stockstalk/internal/change"
	"stockstalk/internal/checkpoint"
	"stockstalk/internal/config"
	"stockstalk/internal/db"
	"stockstalk/internal/graphql"
	"stockstalk/internal/instrument"
	"stockstalk/internal/marketdata"
	"stockstalk/internal/processor"
	"stockstalk/internal/user"
	"stockstalk/internal/watchlist"
	"stockstalk/pkg/logger"
	"syscall"
	"time"
)

func main() {
	// ------------------------------------------------------------
	// Structured logger
	// ------------------------------------------------------------
	structuredLogger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	slog.SetDefault(structuredLogger)

	// ------------------------------------------------------------
	// Configuration
	// ------------------------------------------------------------
	cfg := config.Load()
	appLogger := logger.New()

	// ------------------------------------------------------------
	// Database
	// ------------------------------------------------------------
	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal(err)
	}
	defer database.Close()

	// ------------------------------------------------------------
	// Database migrations
	// ------------------------------------------------------------
	if err := db.Migrate(
		context.Background(),
		database.DB,
		"migrations/001_init.sql",
	); err != nil {
		appLogger.Fatal(err)
	}

	if err := db.Migrate(
		context.Background(),
		database.DB,
		"migrations/002_event_intelligence.sql",
	); err != nil {
		appLogger.Fatal(err)
	}

	if err := db.Migrate(
		context.Background(),
		database.DB,
		"migrations/003_event_context.sql",
	); err != nil {
		appLogger.Fatal(err)
	}

	// ------------------------------------------------------------
	// Application context / graceful shutdown
	// ------------------------------------------------------------
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	// ------------------------------------------------------------
	// Market data source
	// ------------------------------------------------------------
	market, err := marketdata.NewService(cfg.MarketDataFile)
	if err != nil {
		appLogger.Fatal(err)
	}

	marketData, err := market.GetMarketData()
	if err != nil {
		appLogger.Fatal(err)
	}

	// ------------------------------------------------------------
	// Seed instruments
	// ------------------------------------------------------------
	if err := db.SeedInstruments(
		context.Background(),
		database.DB,
		marketData.Instruments,
	); err != nil {
		appLogger.Fatal(err)
	}

	// ------------------------------------------------------------
	// Repositories
	// ------------------------------------------------------------
	userRepo := user.NewRepository(database.DB)

	instrumentRepo := instrument.NewRepository(
		convertInstruments(marketData.Instruments),
	)

	watchlistRepo := watchlist.NewRepository(database.DB)
	checkpointRepo := checkpoint.NewRepository(database.DB)
	changeRepo := change.NewRepository(database.DB)
	marketRepo := marketdata.NewRepository(database.DB)

	// ------------------------------------------------------------
	// Services
	// ------------------------------------------------------------
	userService := user.NewService(userRepo)

	authService := auth.NewService(
		userService,
		cfg.JWTSecret,
	)

	watchlistService := watchlist.NewService(watchlistRepo)
	instrumentService := instrument.NewService(instrumentRepo)
	checkpointService := checkpoint.NewService(checkpointRepo)
	changeService := change.NewService(changeRepo)

	// ------------------------------------------------------------
	// Seed initial market snapshots
	// ------------------------------------------------------------
	for _, quote := range marketData.Quotes {
		if err := marketRepo.SaveSnapshot(ctx, quote); err != nil {
			appLogger.Fatal(err)
		}
	}

	// ------------------------------------------------------------
	// GraphQL resolver
	// ------------------------------------------------------------
	resolver := graphql.NewResolver(
		userService,
		authService,
		watchlistService,
		instrumentService,
		market,
		marketRepo,
		checkpointService,
		changeService,
	)

	// ------------------------------------------------------------
	// GraphQL server
	// ------------------------------------------------------------
	graphqlServer, err := graphql.NewServer(
		resolver,
		authService,
	)
	if err != nil {
		appLogger.Fatal(err)
	}

	// ------------------------------------------------------------
	// Market processor
	// ------------------------------------------------------------
	marketProcessor := processor.New(
		market,
		marketRepo,
		watchlistService,
		changeService,
		func(userID, watchlistID string, event interface{}) {
			if resolver.Broker != nil {
				resolver.Broker.Publish(
					userID,
					watchlistID,
					event,
				)
			}
		},
	)

	// ------------------------------------------------------------
	// Background market worker
	// ------------------------------------------------------------
	marketWorker := processor.NewWorker(
		marketProcessor,
		30*time.Second,
		log.Default(),
	)

	go marketWorker.Start(ctx)

	slog.Info(
		"stockstalk backend initialized",
		"port",
		cfg.Port,
		"market_processor_interval",
		"30s",
	)

	// ------------------------------------------------------------
	// HTTP routes
	// ------------------------------------------------------------
	mux := http.NewServeMux()

	// GraphQL
	mux.Handle(
		"/graphql",
		graphqlServer.Handler(),
	)

	// GraphQL WebSocket subscriptions
	mux.Handle(
		"/graphql/ws",
		graphqlServer.WebSocketHandler(),
	)

	// Metrics
	mux.Handle(
		"/metrics",
		expvar.Handler(),
	)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(
			`{"status":"ok","service":"stockstalk"}`,
		))
	})

	// ------------------------------------------------------------
	// HTTP server
	// ------------------------------------------------------------
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,

		// Protect normal HTTP requests from hanging forever.
		ReadTimeout: 10 * time.Second,

		// IMPORTANT:
		// WebSocket subscriptions are long-lived, so WriteTimeout
		// must remain disabled.
		WriteTimeout: 0,

		IdleTimeout: 60 * time.Second,
	}

	// ------------------------------------------------------------
	// Start HTTP server
	// ------------------------------------------------------------
	go func() {
		slog.Info(
			"server started",
			"address",
			":"+cfg.Port,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			slog.Error(
				"server stopped unexpectedly",
				"error",
				err,
			)

			appLogger.Printf(
				"server stopped: %v",
				err,
			)
		}
	}()

	// ------------------------------------------------------------
	// Wait for shutdown signal
	// ------------------------------------------------------------
	<-ctx.Done()

	slog.Info("shutdown signal received")

	// ------------------------------------------------------------
	// Graceful shutdown
	// ------------------------------------------------------------
	shutdownCtx, stop := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer stop()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error(
			"shutdown error",
			"error",
			err,
		)

		appLogger.Printf(
			"shutdown error: %v",
			err,
		)
	}

	slog.Info("stockstalk backend stopped")
}

// ------------------------------------------------------------
// Convert marketdata instruments to instrument package models.
// ------------------------------------------------------------
func convertInstruments(
	items []marketdata.Instrument,
) []instrument.Instrument {

	result := make(
		[]instrument.Instrument,
		0,
		len(items),
	)

	for _, item := range items {
		result = append(
			result,
			instrument.Instrument{
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
			},
		)
	}

	return result
}
