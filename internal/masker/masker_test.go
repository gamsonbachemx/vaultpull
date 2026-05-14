package masker_test

import (
	"testing"

	"github.com/user/vaultpull/internal/masker"
)

func TestMask_EmptyValue(t *testing.T) {
	m := masker.New()
	if got := m.Mask(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMask_DefaultMask(t *testing.T) {
	m := masker.New()
	got := m.Mask("supersecret")
	if got != masker.DefaultMask {
		t.Errorf("expected %q, got %q", masker.DefaultMask, got)
	}
}

func TestMask_CustomMask(t *testing.T) {
	m := masker.New(masker.WithMask("[REDACTED]"))
	got := m.Mask("supersecret")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestMask_PartialReveal(t *testing.T) {
	m := masker.New(masker.WithPartial(4))
	got := m.Mask("supersecret")
	if got != "supe"+masker.DefaultMask {
		t.Errorf("expected partial mask, got %q", got)
	}
}

func TestMask_PartialReveal_ShortValue(t *testing.T) {
	m := masker.New(masker.WithPartial(4))
	got := m.Mask("abc")
	if got != masker.DefaultMask {
		t.Errorf("short value should be fully masked, got %q", got)
	}
}

func TestMaskMap(t *testing.T) {
	m := masker.New()
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "tok_abc123",
	}
	masked := m.MaskMap(secrets)
	for k, v := range masked {
		if v != masker.DefaultMask {
			t.Errorf("key %s: expected %q, got %q", k, masker.DefaultMask, v)
		}
	}
	if len(masked) != len(secrets) {
		t.Errorf("expected %d entries, got %d", len(secrets), len(masked))
	}
}

func TestMaskMap_DoesNotMutateOriginal(t *testing.T) {
	m := masker.New()
	secrets := map[string]string{"KEY": "realvalue"}
	_ = m.MaskMap(secrets)
	if secrets["KEY"] != "realvalue" {
		t.Error("original map should not be mutated")
	}
}

func TestIsSensitiveKey(t *testing.T) {
	tests := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"API_SECRET", true},
		{"AUTH_TOKEN", true},
		{"PRIVATE_KEY", true},
		{"DATABASE_URL", false},
		{"APP_ENV", false},
		{"PORT", false},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := masker.IsSensitiveKey(tt.key)
			if got != tt.want {
				t.Errorf("IsSensitiveKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
