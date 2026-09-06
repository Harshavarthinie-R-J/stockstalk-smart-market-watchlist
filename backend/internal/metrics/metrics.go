package metrics

import (
	"expvar"
	"sync/atomic"
	"time"
)

var (
	ProcessedSnapshots = expvar.NewInt("stockstalk_processed_snapshots")
	DetectedEvents     = expvar.NewInt("stockstalk_detected_events")
	ProcessingErrors   = expvar.NewInt("stockstalk_processing_errors")
	StaleQuotes        = expvar.NewInt("stockstalk_stale_quotes")
)

var processingLatencyNS int64

func RecordLatency(d time.Duration) {
	atomic.StoreInt64(&processingLatencyNS, d.Nanoseconds())
}

func ProcessingLatency() time.Duration {
	return time.Duration(
		atomic.LoadInt64(&processingLatencyNS),
	)
}