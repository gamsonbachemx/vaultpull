package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

// newTestConfig is a helper that creates a Config pointing at the given test
// server URL and a temporary output file, reducing boilerplate across tests.
func newTestConfig(t *testing.T, serverURL string) *config.Config {
	t.Helper()
	tmpDir := t.TempDir()
	return &config.Config{
		VaultAddress: serverURL,
		VaultToken:   "test-token",
		SecretPath:   "secret/data/app",
		OutputFile:   filepath.Join(tmpDir, ".env"),
	}
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

	cfg := newTestConfig(t, server.URL)

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

	cfg := newTestConfig(t, server.URL)

	syncer, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating syncer: %v", err)
	}

	if err := syncer.Run(); err != nil {
		t.Fatalf("unexpected error running sync: %v", err)
	}

	data, err := os.ReadFile(cfg.OutputFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected output file to have content, got empty file")
	}

	// Verify that each expected key appears in the output file.
	content := string(data)
	for key := range secrets {
		if !strings.Contains(content, key) {
			t.Errorf("expected output file to contain key %q, but it did not", key)
		}
	}
}
