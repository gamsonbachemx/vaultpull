package env

import (
	"strings"
	"testing"
)

func TestParse_EmptyInput(t *testing.T) {
	p := NewParser()
	result, err := p.Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestParse_BasicKeyValue(t *testing.T) {
	input := "FOO=bar\nBAZ=qux\n"
	p := NewParser()
	result, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", result["FOO"])
	}
	if result["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", result["BAZ"])
	}
}

func TestParse_SkipsComments(t *testing.T) {
	input := "# this is a comment\nFOO=bar\n"
	p := NewParser()
	result, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
}

func TestParse_QuotedValues(t *testing.T) {
	input := `DB_URL="postgres://localhost/mydb"` + "\n" +
		`SECRET='my secret value'` + "\n"
	p := NewParser()
	result, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("unexpected DB_URL: %q", result["DB_URL"])
	}
	if result["SECRET"] != "my secret value" {
		t.Errorf("unexpected SECRET: %q", result["SECRET"])
	}
}

func TestParse_ExportPrefix(t *testing.T) {
	input := "export FOO=bar\n"
	p := NewParser()
	result, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", result["FOO"])
	}
}

func TestParse_MissingEquals(t *testing.T) {
	input := "BADLINE\n"
	p := NewParser()
	_, err := p.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func TestParse_ValueWithEquals(t *testing.T) {
	input := "TOKEN=abc=def=ghi\n"
	p := NewParser()
	result, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["TOKEN"] != "abc=def=ghi" {
		t.Errorf("unexpected TOKEN value: %q", result["TOKEN"])
	}
}
