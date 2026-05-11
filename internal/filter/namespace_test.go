package filter_test

import (
	"testing"

	"github.com/yourorg/vaultpull/internal/filter"
)

func TestMatch_NoNamespaces(t *testing.T) {
	f := filter.New(nil)
	if !f.Match("ANY_KEY") {
		t.Error("expected all keys to match when no namespaces configured")
	}
}

func TestMatch_WithNamespaces(t *testing.T) {
	f := filter.New([]string{"APP_", "DB_"})

	cases := []struct {
		key   string
		want  bool
	}{
		{"APP_SECRET", true},
		{"DB_PASSWORD", true},
		{"OTHER_KEY", false},
		{"APP_", true},
		{"", false},
	}

	for _, tc := range cases {
		got := f.Match(tc.key)
		if got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestApply_NoNamespaces(t *testing.T) {
	f := filter.New([]string{})
	secrets := map[string]string{"A": "1", "B": "2"}
	result := f.Apply(secrets)
	if len(result) != len(secrets) {
		t.Errorf("expected %d entries, got %d", len(secrets), len(result))
	}
}

func TestApply_WithNamespaces(t *testing.T) {
	f := filter.New([]string{"APP_"})
	secrets := map[string]string{
		"APP_TOKEN":  "abc",
		"DB_PASS":    "xyz",
		"APP_SECRET": "def",
	}
	result := f.Apply(secrets)
	if len(result) != 2 {
		t.Fatalf("expected 2 filtered entries, got %d", len(result))
	}
	if result["APP_TOKEN"] != "abc" || result["APP_SECRET"] != "def" {
		t.Error("unexpected filtered values")
	}
}

func TestStripPrefix(t *testing.T) {
	f := filter.New([]string{"APP_"})

	cases := []struct {
		input string
		want  string
	}{
		{"APP_TOKEN", "TOKEN"},
		{"DB_PASS", "DB_PASS"},
		{"APP_", ""},
	}

	for _, tc := range cases {
		got := f.StripPrefix(tc.input)
		if got != tc.want {
			t.Errorf("StripPrefix(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
