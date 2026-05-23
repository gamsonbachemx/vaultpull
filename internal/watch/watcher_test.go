package watch_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/vaultpull/internal/watch"
)

type mockSyncer struct {
	calls atomic.Int32
	err   error
}

func (m *mockSyncer) Run(_ context.Context) error {
	m.calls.Add(1)
	return m.err
}

func TestNew_NilSyncer(t *testing.T) {
	_, err := watch.New(nil, watch.Config{Interval: time.Second})
	if err == nil {
		t.Fatal("expected error for nil syncer")
	}
}

func TestNew_ZeroInterval(t *testing.T) {
	_, err := watch.New(&mockSyncer{}, watch.Config{Interval: 0})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNew_Valid(t *testing.T) {
	w, err := watch.New(&mockSyncer{}, watch.Config{Interval: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestRun_ImmediateSyncOnStart(t *testing.T) {
	s := &mockSyncer{}
	w, _ := watch.New(s, watch.Config{Interval: 10 * time.Second})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)

	if s.calls.Load() < 1 {
		t.Error("expected at least one immediate sync call")
	}
}

func TestRun_TicksRepeatedly(t *testing.T) {
	s := &mockSyncer{}
	w, _ := watch.New(s, watch.Config{Interval: 30 * time.Millisecond})

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)

	if s.calls.Load() < 3 {
		t.Errorf("expected >= 3 calls, got %d", s.calls.Load())
	}
}

func TestRun_CallsOnErrorOnFailure(t *testing.T) {
	s := &mockSyncer{err: errors.New("vault unavailable")}
	var caught error
	w, _ := watch.New(s, watch.Config{
		Interval: 10 * time.Second,
		OnError: func(err error) { caught = err },
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)

	if caught == nil {
		t.Error("expected OnError to be called")
	}
}

func TestRun_ReturnsContextError(t *testing.T) {
	s := &mockSyncer{}
	w, _ := watch.New(s, watch.Config{Interval: time.Second})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := w.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestRun_NoSyncWhenContextAlreadyCanceled(t *testing.T) {
	s := &mockSyncer{}
	w, _ := watch.New(s, watch.Config{Interval: time.Second})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_ = w.Run(ctx)

	if s.calls.Load() != 0 {
		t.Errorf("expected no sync calls when context is pre-cancelled, got %d", s.calls.Load())
	}
}
