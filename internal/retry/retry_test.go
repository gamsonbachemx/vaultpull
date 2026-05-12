package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"vaultpull/internal/retry"
)

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	sentinel := errors.New("transient error")
	cfg := retry.Config{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}

	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success on third attempt, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_AllAttemptsFail(t *testing.T) {
	sentinel := errors.New("permanent error")
	cfg := retry.Config{MaxAttempts: 2, Delay: time.Millisecond, Multiplier: 1.0}

	err := retry.Do(context.Background(), cfg, func() error {
		return sentinel
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel in error chain, got %v", err)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := retry.Config{MaxAttempts: 5, Delay: time.Millisecond, Multiplier: 1.0}
	err := retry.Do(ctx, cfg, func() error {
		return errors.New("should not matter")
	})
	if err == nil {
		t.Fatal("expected error due to cancelled context")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := retry.DefaultConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
}
