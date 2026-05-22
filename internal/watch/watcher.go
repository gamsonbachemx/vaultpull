// Package watch provides periodic re-sync of Vault secrets into local .env files.
package watch

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Syncer is the interface that the watcher calls on each tick.
type Syncer interface {
	Run(ctx context.Context) error
}

// Config holds configuration for the watcher.
type Config struct {
	// Interval between sync attempts.
	Interval time.Duration
	// OnError is called when a sync fails. If nil, errors are logged to stderr.
	OnError func(err error)
}

// Watcher runs a Syncer on a fixed interval until the context is cancelled.
type Watcher struct {
	cfg    Config
	syncer Syncer
}

// New creates a new Watcher. Returns an error if syncer is nil or interval is zero.
func New(syncer Syncer, cfg Config) (*Watcher, error) {
	if syncer == nil {
		return nil, fmt.Errorf("watch: syncer must not be nil")
	}
	if cfg.Interval <= 0 {
		return nil, fmt.Errorf("watch: interval must be positive")
	}
	if cfg.OnError == nil {
		cfg.OnError = func(err error) {
			log.Printf("watch: sync error: %v", err)
		}
	}
	return &Watcher{cfg: cfg, syncer: syncer}, nil
}

// Run starts the watch loop. It performs an immediate sync, then repeats on
// every tick. Blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	if err := w.tick(ctx); err != nil {
		w.cfg.OnError(err)
	}

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(ctx); err != nil {
				w.cfg.OnError(err)
			}
		}
	}
}

func (w *Watcher) tick(ctx context.Context) error {
	return w.syncer.Run(ctx)
}
