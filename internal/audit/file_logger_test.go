package audit_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"vaultpull/internal/audit"
)

func TestNewFileLogger_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logs", "audit.log")

	fl, err := audit.NewFileLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer fl.Close()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected log file to be created")
	}
}

func TestNewFileLogger_WritesAndClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	fl, err := audit.NewFileLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := fl.Sync("/secret/app", 3); err != nil {
		t.Fatalf("log error: %v", err)
	}

	if err := fl.Close(); err != nil {
		t.Fatalf("close error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(data[:len(data)-1], &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Type != audit.EventSync {
		t.Errorf("expected EventSync, got %q", e.Type)
	}
}

func TestDiscard_NoError(t *testing.T) {
	l := audit.Discard()
	if err := l.Log(audit.EventRead, "ignored", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewFileLogger_InvalidPath(t *testing.T) {
	// Use a path where a file exists where a directory is expected.
	dir := t.TempDir()
	blocking := filepath.Join(dir, "blocked")
	if err := os.WriteFile(blocking, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := audit.NewFileLogger(filepath.Join(blocking, "audit.log"))
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
