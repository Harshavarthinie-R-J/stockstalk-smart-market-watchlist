package marketdata

import (
	"context"
	"time"
)

func Retry(
	ctx context.Context,
	attempts int,
	initialDelay time.Duration,
	fn func() error,
) error {

	if attempts <= 0 {
		attempts = 3
	}

	if initialDelay <= 0 {
		initialDelay = 250 * time.Millisecond
	}

	var err error
	delay := initialDelay

	for attempt := 0; attempt < attempts; attempt++ {

		if err = fn(); err == nil {
			return nil
		}

		if attempt == attempts-1 {
			break
		}

		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()

		case <-timer.C:
		}

		delay *= 2

		if delay > 5*time.Second {
			delay = 5 * time.Second
		}
	}

	return err
}
