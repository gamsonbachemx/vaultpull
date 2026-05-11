package config_test

import (
	"os"
	"testing"

	"github.com/yourorg/vaultpull/internal/config"
)

func TestLoad_MissingToken(t *testing.T) {
	os.Setenv("VAULT_ADDR", "http://localhost:8200")
	os.Unsetenv("VAULT_TOKEN")
	defer os.Unsetenv("VAULT_ADDR")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for missing VAULT_TOKEN")
	}
}

func TestLoad_Defaults(t *testing.T) {
	os.Setenv("VAULT_ADDR", "http://localhost:8200")
	os.Setenv("VAULT_TOKEN", "test-token")
	defer os.Unsetenv("VAULT_ADDR")
	defer os.Unsetenv("VAULT_TOKEN")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SecretPath != "secret/data/app" {
		t.Errorf("expected default secret path, got %q", cfg.SecretPath)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default output file, got %q", cfg.OutputFile)
	}
	if cfg.StripPrefix {
		t.Error("expected StripPrefix to default to false")
	}
}

func TestLoad_CustomValues(t *testing.T) {
	os.Setenv("VAULT_ADDR", "http://vault:8200")
	os.Setenv("VAULT_TOKEN", "s.abc123")
	os.Setenv("VAULT_SECRET_PATH", "secret/data/myapp")
	os.Setenv("OUTPUT_FILE", ".env.local")
	os.Setenv("VAULT_NAMESPACES", "APP_, DB_")
	os.Setenv("STRIP_PREFIX", "true")
	defer func() {
		for _, k := range []string{"VAULT_ADDR", "VAULT_TOKEN", "VAULT_SECRET_PATH", "OUTPUT_FILE", "VAULT_NAMESPACES", "STRIP_PREFIX"} {
			os.Unsetenv(k)
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Namespaces) != 2 {
		t.Fatalf("expected 2 namespaces, got %d", len(cfg.Namespaces))
	}
	if cfg.Namespaces[0] != "APP_" || cfg.Namespaces[1] != "DB_" {
		t.Errorf("unexpected namespaces: %v", cfg.Namespaces)
	}
	if !cfg.StripPrefix {
		t.Error("expected StripPrefix to be true")
	}
}
