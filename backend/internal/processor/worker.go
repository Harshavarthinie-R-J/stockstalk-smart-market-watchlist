package processor

import (
	"context"
	"log"
	"log/slog"
	"time"
)

type Worker struct {
	Processor *Processor
	Interval  time.Duration
	Logger    *log.Logger
}

func NewWorker(
	processor *Processor,
	interval time.Duration,
	logger *log.Logger,
) *Worker {

	if interval <= 0 {
		interval = 30 * time.Second
	}

	if logger == nil {
		logger = log.Default()
	}

	return &Worker{
		Processor: processor,
		Interval:  interval,
		Logger:    logger,
	}
}

func (w *Worker) Start(ctx context.Context) {
	if w.Processor == nil {
		slog.Error("market processor worker has no processor")
		return
	}

	slog.Info(
		"market processor started",
		"interval",
		w.Interval.String(),
	)

	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	// Run immediately when the application starts.
	w.run(ctx)

	for {
		select {
		case <-ctx.Done():
			slog.Info("market processor stopped")
			return

		case <-ticker.C:
			w.run(ctx)
		}
	}
}

func (w *Worker) run(ctx context.Context) {
	started := time.Now()

	if err := w.Processor.Process(ctx); err != nil {
		slog.Error(
			"market processor error",
			"error",
			err,
			"duration",
			time.Since(started).String(),
		)

		if w.Logger != nil {
			w.Logger.Printf(
				"market processor error: %v",
				err,
			)
		}

		return
	}

	slog.Info(
		"market processor completed",
		"duration",
		time.Since(started).String(),
	)
}
