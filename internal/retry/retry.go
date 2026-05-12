package retry

import (
	"context"
	"fmt"
	"time"
)

// Config holds retry configuration.
type Config struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Do executes fn up to MaxAttempts times, backing off between attempts.
// It returns the last error if all attempts fail.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.Multiplier <= 0 {
		cfg.Multiplier = 1.0
	}

	delay := cfg.Delay
	var lastErr error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry aborted: %w", err)
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("retry aborted: %w", ctx.Err())
		case <-time.After(delay):
		}

		delay = time.Duration(float64(delay) * cfg.Multiplier)
	}

	return fmt.Errorf("all %d attempts failed: %w", cfg.MaxAttempts, lastErr)
}
