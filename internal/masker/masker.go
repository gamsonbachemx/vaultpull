// Package masker provides utilities for masking sensitive secret values
// before displaying them in logs, diffs, or terminal output.
package masker

import "strings"

const (
	// DefaultMask is the string used to replace sensitive values.
	DefaultMask = "***"
	// RevealChars is the number of characters to reveal at the start of a value.
	RevealChars = 4
)

// Masker masks secret values for safe display.
type Masker struct {
	mask        string
	revealChars int
	partial     bool
}

// Option configures a Masker.
type Option func(*Masker)

// WithMask sets a custom mask string.
func WithMask(mask string) Option {
	return func(m *Masker) {
		m.mask = mask
	}
}

// WithPartial enables partial reveal mode, showing the first N characters.
func WithPartial(chars int) Option {
	return func(m *Masker) {
		m.partial = true
		m.revealChars = chars
	}
}

// New creates a new Masker with the given options.
func New(opts ...Option) *Masker {
	m := &Masker{
		mask:        DefaultMask,
		revealChars: RevealChars,
		partial:     false,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Mask replaces a secret value with the mask string.
func (m *Masker) Mask(value string) string {
	if value == "" {
		return ""
	}
	if m.partial && len(value) > m.revealChars {
		return value[:m.revealChars] + m.mask
	}
	return m.mask
}

// MaskMap masks all values in a map, returning a new map safe for display.
func (m *Masker) MaskMap(secrets map[string]string) map[string]string {
	masked := make(map[string]string, len(secrets))
	for k, v := range secrets {
		masked[k] = m.Mask(v)
	}
	return masked
}

// IsSensitiveKey reports whether the key name suggests a sensitive value.
func IsSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	sensitivePatterns := []string{"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL", "AUTH"}
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
