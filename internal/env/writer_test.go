package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_EmptySecrets(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".env")

	w := NewWriter(path)
	if err := w.Write(map[string]string{}); err != nil {
		t.Fatalf("expected no error for empty secrets, got: %v", err)
	}

	// File should not be created for empty secrets.
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected file to not exist for empty secrets")
	}
}

func TestWrite_BasicSecrets(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".env")

	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
	}

	w := NewWriter(path)
	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "API_KEY=abc123") {
		t.Errorf("expected API_KEY in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PASSWORD=s3cr3t") {
		t.Errorf("expected DB_PASSWORD in output, got:\n%s", content)
	}

	// Verify sorted order: API_KEY should come before DB_PASSWORD.
	apiIdx := strings.Index(content, "API_KEY")
	dbIdx := strings.Index(content, "DB_PASSWORD")
	if apiIdx > dbIdx {
		t.Errorf("expected keys to be sorted alphabetically")
	}
}

func TestWrite_QuotesSpecialValues(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".env")

	secrets := map[string]string{
		"MY_VAR": "hello world",
	}

	w := NewWriter(path)
	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), `"hello world"`) {
		t.Errorf("expected value with space to be quoted, got: %s", string(data))
	}
}

func TestFormatKey(t *testing.T) {
	cases := []struct {
		namespace, field, expected string
	}{
		{"myapp/database", "password", "MYAPP_DATABASE_PASSWORD"},
		{"my-service", "api-key", "MY_SERVICE_API_KEY"},
		{"infra.prod", "token", "INFRA_PROD_TOKEN"},
	}

	for _, tc := range cases {
		got := FormatKey(tc.namespace, tc.field)
		if got != tc.expected {
			t.Errorf("FormatKey(%q, %q) = %q, want %q", tc.namespace, tc.field, got, tc.expected)
		}
	}
}
