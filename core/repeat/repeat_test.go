package repeat

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExecSuccess(t *testing.T) {
	op := func(ctx context.Context) error {
		return nil // Operation succeeds
	}

	if err := Exec(context.Background(), op); err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestExecFailure(t *testing.T) {
	op := func(ctx context.Context) error {
		return errors.New("failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // Increased timeout
	defer cancel()

	if err := Exec(ctx, op, WithMaxRetries(1), WithMinTimeWait(time.Millisecond), WithMaxTimeWait(time.Millisecond*10)); err == nil || err.Error() != "failed" {
		t.Fatalf("Expected 'failed' error, got: %v", err)
	}
}

func TestExecInfiniteRetries(t *testing.T) {
	const maxCount = 100
	count := 0
	op := func(ctx context.Context) error {
		count++
		if count >= maxCount { // Succeeds after maxCount retries
			return nil
		}
		return errors.New("failed")
	}

	// Increasing the context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Decreasing the min and max wait time
	if err := Exec(ctx, op, WithMinTimeWait(10*time.Millisecond), WithMaxTimeWait(50*time.Millisecond)); err != nil { // Infinite retries by default
		t.Fatalf("Expected no error, got: %v", err)
	}

	if count != maxCount {
		t.Fatalf("Expected %d executions, got: %d", maxCount, count)
	}
}

func TestExecJitter(t *testing.T) {
	count := 0
	op := func(ctx context.Context) error {
		count++
		return errors.New("failed")
	}

	// Increasing the context timeout to allow more time for retries
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	startTime := time.Now()
	if err := Exec(ctx, op, WithMaxRetries(1), WithMinTimeWait(10*time.Millisecond), WithMaxTimeWait(50*time.Millisecond)); err == nil || err.Error() != "failed" {
		t.Fatalf("Expected 'failed' error, got: %v", err)
	}

	duration := time.Since(startTime)
	if duration < 10*time.Millisecond {
		t.Fatalf("Expected duration to be at least 10 milliseconds, got: %v", duration)
	}
}

func TestExecMinAndMaxTimeouts(t *testing.T) {
	op := func(ctx context.Context) error {
		return errors.New("failed")
	}

	minTimeout := 50 * time.Millisecond
	maxTimeout := 100 * time.Millisecond

	// Context with enough time to allow for retries
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Function to measure retry time
	retryTime := func() time.Duration {
		startTime := time.Now()
		Exec(ctx, op, WithMaxRetries(1), WithMinTimeWait(minTimeout), WithMaxTimeWait(maxTimeout))
		return time.Since(startTime)
	}

	for i := 0; i < 5; i++ { // Testing multiple retries to check the range
		duration := retryTime()
		if duration < minTimeout || duration > maxTimeout {
			t.Fatalf("Expected duration to be within %v and %v, got: %v", minTimeout, maxTimeout, duration)
		}
	}
}
