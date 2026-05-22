package env

import (
	"testing"
)

func TestSanitizeKey_NoNormalize(t *testing.T) {
	s := NewSanitizer()
	got := s.SanitizeKey("my-key")
	if got != "my-key" {
		t.Errorf("expected 'my-key', got %q", got)
	}
}

func TestSanitizeKey_Normalize(t *testing.T) {
	s := NewSanitizer(WithNormalizeKeys())
	cases := []struct {
		input string
		want  string
	}{
		{"my-key", "MY_KEY"},
		{"db.host", "DB_HOST"},
		{"  leading", "LEADING"},
		{"ALREADY_VALID", "ALREADY_VALID"},
		{"mixed-Case_Key", "MIXED_CASE_KEY"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := s.SanitizeKey(tc.input)
			if got != tc.want {
				t.Errorf("SanitizeKey(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestSanitizeKey_EmptyAfterNormalize(t *testing.T) {
	s := NewSanitizer(WithNormalizeKeys())
	got := s.SanitizeKey("---")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestSanitizeValue_NoStrip(t *testing.T) {
	s := NewSanitizer()
	input := "hello\x01world"
	got := s.SanitizeValue(input)
	if got != input {
		t.Errorf("expected value unchanged, got %q", got)
	}
}

func TestSanitizeValue_StripControlChars(t *testing.T) {
	s := NewSanitizer(WithStripControlChars())
	input := "hello\x01\x02world"
	want := "helloworld"
	got := s.SanitizeValue(input)
	if got != want {
		t.Errorf("SanitizeValue(%q) = %q, want %q", input, got, want)
	}
}

func TestSanitizeValue_PreservesNewlineAndTab(t *testing.T) {
	s := NewSanitizer(WithStripControlChars())
	input := "line1\nline2\ttabbed"
	got := s.SanitizeValue(input)
	if got != input {
		t.Errorf("expected newline/tab preserved, got %q", got)
	}
}

func TestSanitizeAll_SkipsEmptyKeys(t *testing.T) {
	s := NewSanitizer(WithNormalizeKeys())
	secrets := map[string]string{
		"---":    "should be dropped",
		"DB_URL": "postgres://localhost",
	}
	out := s.SanitizeAll(secrets)
	if _, ok := out[""]; ok {
		t.Error("empty key should not appear in output")
	}
	if out["DB_URL"] != "postgres://localhost" {
		t.Errorf("expected DB_URL to be preserved")
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}

func TestSanitizeAll_NormalizesAndStrips(t *testing.T) {
	s := NewSanitizer(WithNormalizeKeys(), WithStripControlChars())
	secrets := map[string]string{
		"api.key": "secret\x00value",
	}
	out := s.SanitizeAll(secrets)
	if v, ok := out["API_KEY"]; !ok {
		t.Error("expected API_KEY in output")
	} else if v != "secretvalue" {
		t.Errorf("expected 'secretvalue', got %q", v)
	}
}
