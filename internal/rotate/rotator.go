// Package rotate provides functionality to detect and rotate secrets
// that have been synced into local .env files.
package rotate

import (
	"context"
	"fmt"
	"time"
)

// Secret represents a secret entry with metadata.
type Secret struct {
	Key       string
	Value     string
	CreatedAt time.Time
	TTL       time.Duration
}

// IsExpired reports whether the secret has exceeded its TTL.
func (s Secret) IsExpired() bool {
	if s.TTL == 0 {
		return false
	}
	return time.Since(s.CreatedAt) > s.TTL
}

// VaultWriter is the interface for writing new secret values back to Vault.
type VaultWriter interface {
	WriteSecret(ctx context.Context, path, key, value string) error
}

// Generator produces a new secret value for a given key.
type Generator func(key string) (string, error)

// Config holds rotation configuration.
type Config struct {
	Paths     []string
	TTL       time.Duration
	Generator Generator
}

// Rotator checks secrets for expiry and rotates them via Vault.
type Rotator struct {
	cfg    *Config
	vault  VaultWriter
}

// New creates a new Rotator. Returns an error if cfg or vault is nil.
func New(cfg *Config, vault VaultWriter) (*Rotator, error) {
	if cfg == nil {
		return nil, fmt.Errorf("rotate: config must not be nil")
	}
	if vault == nil {
		return nil, fmt.Errorf("rotate: vault writer must not be nil")
	}
	if cfg.Generator == nil {
		return nil, fmt.Errorf("rotate: generator must not be nil")
	}
	return &Rotator{cfg: cfg, vault: vault}, nil
}

// Rotate iterates over secrets and writes a new value for any that are expired.
// Returns the number of secrets rotated and any error encountered.
func (r *Rotator) Rotate(ctx context.Context, secrets []Secret) (int, error) {
	rotated := 0
	for _, s := range secrets {
		if !s.IsExpired() {
			continue
		}
		newVal, err := r.cfg.Generator(s.Key)
		if err != nil {
			return rotated, fmt.Errorf("rotate: generate value for %q: %w", s.Key, err)
		}
		for _, path := range r.cfg.Paths {
			if err := r.vault.WriteSecret(ctx, path, s.Key, newVal); err != nil {
				return rotated, fmt.Errorf("rotate: write secret %q at %q: %w", s.Key, path, err)
			}
		}
		rotated++
	}
	return rotated, nil
}
