package rotate_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/rotate"
)

// stubVaultWriter records calls to WriteSecret.
type stubVaultWriter struct {
	calls []struct{ path, key, value string }
	err   error
}

func (s *stubVaultWriter) WriteSecret(_ context.Context, path, key, value string) error {
	s.calls = append(s.calls, struct{ path, key, value string }{path, key, value})
	return s.err
}

func newSecret(key string, age, ttl time.Duration) rotate.Secret {
	return rotate.Secret{
		Key:       key,
		Value:     "oldval",
		CreatedAt: time.Now().Add(-age),
		TTL:       ttl,
	}
}

func TestNew_NilConfig(t *testing.T) {
	_, err := rotate.New(nil, &stubVaultWriter{})
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestNew_NilVault(t *testing.T) {
	cfg := &rotate.Config{Generator: func(string) (string, error) { return "", nil }}
	_, err := rotate.New(cfg, nil)
	if err == nil {
		t.Fatal("expected error for nil vault writer")
	}
}

func TestNew_NilGenerator(t *testing.T) {
	_, err := rotate.New(&rotate.Config{}, &stubVaultWriter{})
	if err == nil {
		t.Fatal("expected error for nil generator")
	}
}

func TestRotate_SkipsNonExpired(t *testing.T) {
	vault := &stubVaultWriter{}
	cfg := &rotate.Config{
		Paths:     []string{"secret/app"},
		Generator: func(string) (string, error) { return "newval", nil },
	}
	r, _ := rotate.New(cfg, vault)
	secrets := []rotate.Secret{newSecret("TOKEN", 1*time.Minute, 1*time.Hour)}
	n, err := r.Rotate(context.Background(), secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 rotated, got %d", n)
	}
	if len(vault.calls) != 0 {
		t.Errorf("expected no writes, got %d", len(vault.calls))
	}
}

func TestRotate_RotatesExpired(t *testing.T) {
	vault := &stubVaultWriter{}
	cfg := &rotate.Config{
		Paths:     []string{"secret/app"},
		Generator: func(string) (string, error) { return "newval", nil },
	}
	r, _ := rotate.New(cfg, vault)
	secrets := []rotate.Secret{newSecret("TOKEN", 2*time.Hour, 1*time.Hour)}
	n, err := r.Rotate(context.Background(), secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 rotated, got %d", n)
	}
	if len(vault.calls) != 1 || vault.calls[0].key != "TOKEN" || vault.calls[0].value != "newval" {
		t.Errorf("unexpected vault calls: %+v", vault.calls)
	}
}

func TestRotate_GeneratorError(t *testing.T) {
	vault := &stubVaultWriter{}
	cfg := &rotate.Config{
		Paths:     []string{"secret/app"},
		Generator: func(string) (string, error) { return "", errors.New("gen failed") },
	}
	r, _ := rotate.New(cfg, vault)
	secrets := []rotate.Secret{newSecret("TOKEN", 2*time.Hour, 1*time.Hour)}
	_, err := r.Rotate(context.Background(), secrets)
	if err == nil {
		t.Fatal("expected error from generator")
	}
}

func TestSecret_IsExpired_NoTTL(t *testing.T) {
	s := rotate.Secret{Key: "K", CreatedAt: time.Now().Add(-24 * time.Hour)}
	if s.IsExpired() {
		t.Error("secret with no TTL should never be expired")
	}
}
