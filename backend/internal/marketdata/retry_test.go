package marketdata

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryEventuallySucceeds(t *testing.T) {
	attempts := 0

	err := Retry(
		context.Background(),
		3,
		time.Millisecond,
		func() error {
			attempts++

			if attempts < 3 {
				return errors.New("temporary failure")
			}

			return nil
		},
	)

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryReturnsFinalError(t *testing.T) {
	attempts := 0

	expected := errors.New("permanent failure")

	err := Retry(
		context.Background(),
		3,
		time.Millisecond,
		func() error {
			attempts++
			return expected
		},
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Retry(
		ctx,
		3,
		time.Millisecond,
		func() error {
			return errors.New("temporary failure")
		},
	)

	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
