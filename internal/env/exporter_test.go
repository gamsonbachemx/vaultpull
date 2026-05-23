package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestExport_EmptySecrets(t *testing.T) {
	var buf bytes.Buffer
	e := NewExporter(&buf)
	if err := e.Export(map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestExport_BasicSecrets(t *testing.T) {
	var buf bytes.Buffer
	e := NewExporter(&buf)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	if err := e.Export(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export DB_HOST=localhost\n") {
		t.Errorf("missing DB_HOST line in:\n%s", out)
	}
	if !strings.Contains(out, "export DB_PORT=5432\n") {
		t.Errorf("missing DB_PORT line in:\n%s", out)
	}
}

func TestExport_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	e := NewExporter(&buf)
	secrets := map[string]string{
		"Z_KEY": "z",
		"A_KEY": "a",
		"M_KEY": "m",
	}
	if err := e.Export(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "export A_KEY") {
		t.Errorf("expected A_KEY first, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "export Z_KEY") {
		t.Errorf("expected Z_KEY last, got %q", lines[2])
	}
}

func TestExport_QuotesSpecialValues(t *testing.T) {
	var buf bytes.Buffer
	e := NewExporter(&buf)
	secrets := map[string]string{
		"API_KEY": "p@ssw0rd!",
		"GREETING": "hello world",
	}
	if err := e.Export(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export API_KEY='p@ssw0rd!'\n") {
		t.Errorf("expected quoted API_KEY in:\n%s", out)
	}
	if !strings.Contains(out, "export GREETING='hello world'\n") {
		t.Errorf("expected quoted GREETING in:\n%s", out)
	}
}

func TestExport_EscapesSingleQuotes(t *testing.T) {
	var buf bytes.Buffer
	e := NewExporter(&buf)
	secrets := map[string]string{
		"MSG": "it's alive",
	}
	if err := e.Export(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	expected := "export MSG='it'\\''s alive'\n"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestEscapeShellValue_NoQuotingNeeded(t *testing.T) {
	cases := []string{"simple", "with-dash", "with_under", "path/value", "host:port"}
	for _, v := range cases {
		got := escapeShellValue(v)
		if got != v {
			t.Errorf("escapeShellValue(%q) = %q, want %q", v, got, v)
		}
	}
}
