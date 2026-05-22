package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeEnvFile: %v", err)
	}
	return p
}

func TestRead_FileNotExist(t *testing.T) {
	r := NewReader("/nonexistent/.env")
	got, err := r.Read()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestRead_BasicKeyValue(t *testing.T) {
	dir := t.TempDir()
	p := writeEnvFile(t, dir, ".env", "FOO=bar\nBAZ=qux\n")

	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" {
		t.Errorf("FOO: want bar, got %q", got["FOO"])
	}
	if got["BAZ"] != "qux" {
		t.Errorf("BAZ: want qux, got %q", got["BAZ"])
	}
}

func TestRead_SkipsCommentsAndBlanks(t *testing.T) {
	dir := t.TempDir()
	p := writeEnvFile(t, dir, ".env", "# comment\n\nKEY=value\n")

	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got))
	}
	if got["KEY"] != "value" {
		t.Errorf("KEY: want value, got %q", got["KEY"])
	}
}

func TestRead_QuotedValues(t *testing.T) {
	dir := t.TempDir()
	p := writeEnvFile(t, dir, ".env", `SINGLE='hello world'
DOUBLE="hello world"
`)

	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["SINGLE"] != "hello world" {
		t.Errorf("SINGLE: want 'hello world', got %q", got["SINGLE"])
	}
	if got["DOUBLE"] != "hello world" {
		t.Errorf("DOUBLE: want 'hello world', got %q", got["DOUBLE"])
	}
}

func TestRead_ExportPrefix(t *testing.T) {
	dir := t.TempDir()
	p := writeEnvFile(t, dir, ".env", "export TOKEN=abc123\n")

	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["TOKEN"] != "abc123" {
		t.Errorf("TOKEN: want abc123, got %q", got["TOKEN"])
	}
}

func TestRead_MissingSeparator(t *testing.T) {
	dir := t.TempDir()
	p := writeEnvFile(t, dir, ".env", "BADLINE\n")

	r := NewReader(p)
	_, err := r.Read()
	if err == nil {
		t.Fatal("expected error for line missing '=', got nil")
	}
}
