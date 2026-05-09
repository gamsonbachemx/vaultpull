package config

import (
	"os"
	"testing"
)

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when VAULT_TOKEN is missing, got nil")
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "test-token")
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULTPULL_NAMESPACE")
	os.Unsetenv("VAULTPULL_OUTPUT")
	os.Unsetenv("VAULTPULL_PREFIX")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected default VaultAddr, got %q", cfg.VaultAddr)
	}
	if cfg.Namespace != "secret" {
		t.Errorf("expected default Namespace 'secret', got %q", cfg.Namespace)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default OutputFile '.env', got %q", cfg.OutputFile)
	}
	if cfg.Prefix != "" {
		t.Errorf("expected empty Prefix, got %q", cfg.Prefix)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "s.abc123")
	t.Setenv("VAULT_ADDR", "https://vault.example.com")
	t.Setenv("VAULTPULL_NAMESPACE", "myapp/prod")
	t.Setenv("VAULTPULL_OUTPUT", ".env.prod")
	t.Setenv("VAULTPULL_PREFIX", "myapp")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultAddr != "https://vault.example.com" {
		t.Errorf("unexpected VaultAddr: %q", cfg.VaultAddr)
	}
	if cfg.Namespace != "myapp/prod" {
		t.Errorf("unexpected Namespace: %q", cfg.Namespace)
	}
	if cfg.OutputFile != ".env.prod" {
		t.Errorf("unexpected OutputFile: %q", cfg.OutputFile)
	}
	if cfg.Prefix != "MYAPP" {
		t.Errorf("expected prefix to be uppercased 'MYAPP', got %q", cfg.Prefix)
	}
}
