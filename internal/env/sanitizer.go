package env

import (
	"regexp"
	"strings"
	"unicode"
)

// Sanitizer cleans and normalizes environment variable keys and values
// before writing to .env files.
type Sanitizer struct {
	stripControlChars bool
	normalizeKeys     bool
	invalidKeyChars   *regexp.Regexp
}

// SanitizerOption configures a Sanitizer.
type SanitizerOption func(*Sanitizer)

// WithNormalizeKeys enables uppercasing and replacing invalid characters in keys.
func WithNormalizeKeys() SanitizerOption {
	return func(s *Sanitizer) {
		s.normalizeKeys = true
	}
}

// WithStripControlChars enables stripping of control characters from values.
func WithStripControlChars() SanitizerOption {
	return func(s *Sanitizer) {
		s.stripControlChars = true
	}
}

// NewSanitizer creates a Sanitizer with the given options.
func NewSanitizer(opts ...SanitizerOption) *Sanitizer {
	s := &Sanitizer{
		invalidKeyChars: regexp.MustCompile(`[^A-Z0-9_]`),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// SanitizeKey cleans an environment variable key.
// If normalizeKeys is enabled, it uppercases and replaces invalid characters with underscores.
func (s *Sanitizer) SanitizeKey(key string) string {
	if s.normalizeKeys {
		key = strings.ToUpper(key)
		key = s.invalidKeyChars.ReplaceAllString(key, "_")
		key = strings.Trim(key, "_")
	}
	return key
}

// SanitizeValue cleans an environment variable value.
// If stripControlChars is enabled, control characters (except newlines in quoted context) are removed.
func (s *Sanitizer) SanitizeValue(value string) string {
	if !s.stripControlChars {
		return value
	}
	var b strings.Builder
	for _, r := range value {
		if r == '\n' || r == '\t' || !unicode.IsControl(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// SanitizeAll applies SanitizeKey and SanitizeValue to an entire secrets map,
// returning a new map with sanitized keys and values.
func (s *Sanitizer) SanitizeAll(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		cleanKey := s.SanitizeKey(k)
		if cleanKey == "" {
			continue
		}
		out[cleanKey] = s.SanitizeValue(v)
	}
	return out
}
