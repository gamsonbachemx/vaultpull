package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"vaultpull/internal/config"
)

func newVaultTestServer(t *testing.T, secrets map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"data": secrets,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
}

func TestNew_NilConfig(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil config, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	server := newVaultTestServer(t, map[string]interface{}{})
	defer server.Close()

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, ".env")

	cfg := &config.Config{
		VaultAddress: server.URL,
		VaultToken:   "test-token",
		SecretPath:   "secret/data/app",
		OutputFile:   outFile,
	}

	syncer, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if syncer == nil {
		t.Fatal("expected non-nil syncer")
	}
}

func TestRun_WritesSecrets(t *testing.T) {
	secrets := map[string]interface{}{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	server := newVaultTestServer(t, secrets)
	defer server.Close()

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, ".env")

	cfg := &config.Config{
		VaultAddress: server.URL,
		VaultToken:   "test-token",
		SecretPath:   "secret/data/app",
		OutputFile:   outFile,
	}

	syncer, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating syncer: %v", err)
	}

	if err := syncer.Run(); err != nil {
		t.Fatalf("unexpected error running sync: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected output file to have content, got empty file")
	}
}
