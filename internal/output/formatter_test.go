package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	f := New(FormatDotenv, nil)
	if f.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWrite_DotenvFormat(t *testing.T) {
	var buf bytes.Buffer
	f := New(FormatDotenv, &buf)

	err := f.Write(map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", got)
	}
}

func TestWrite_ExportFormat(t *testing.T) {
	var buf bytes.Buffer
	f := New(FormatExport, &buf)

	err := f.Write(map[string]string{"DB_URL": "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "export DB_URL=localhost") {
		t.Errorf("expected export prefix in output, got: %s", got)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	f := New(FormatJSON, &buf)

	err := f.Write(map[string]string{"API_KEY": "secret123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "{") || !strings.Contains(got, "API_KEY") {
		t.Errorf("expected JSON output with API_KEY, got: %s", got)
	}
}

func TestWrite_EmptySecrets(t *testing.T) {
	var buf bytes.Buffer
	f := New(FormatDotenv, &buf)

	err := f.Write(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty secrets")
	}
}

func TestQuoteIfNeeded(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with space", `"with space"`},
		{"has#hash", `"has#hash"`},
		{"has$dollar", `"has$dollar"`},
		{"normal123", "normal123"},
	}

	for _, tt := range tests {
		got := quoteIfNeeded(tt.input)
		if got != tt.expected {
			t.Errorf("quoteIfNeeded(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
