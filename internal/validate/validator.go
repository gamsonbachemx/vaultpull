// Package validate provides secret key and value validation for vaultpull.
package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// validKeyPattern matches valid environment variable key names.
var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Result holds the outcome of a validation check.
type Result struct {
	Key     string
	Warning string
	Err     error
}

// IsValid returns true if the result has no error.
func (r Result) IsValid() bool {
	return r.Err == nil
}

// Validator checks secrets for key naming and value constraints.
type Validator struct {
	maxValueLen int
	warnOnLower bool
}

// Option configures a Validator.
type Option func(*Validator)

// WithMaxValueLen sets the maximum allowed value length.
func WithMaxValueLen(n int) Option {
	return func(v *Validator) { v.maxValueLen = n }
}

// WithWarnOnLowercase emits a warning when keys contain lowercase letters.
func WithWarnOnLowercase() Option {
	return func(v *Validator) { v.warnOnLower = true }
}

// New creates a Validator with the provided options.
func New(opts ...Option) *Validator {
	v := &Validator{maxValueLen: 65536}
	for _, o := range opts {
		o(v)
	}
	return v
}

// ValidateKey checks whether a single key is a legal env var name.
func (v *Validator) ValidateKey(key string) Result {
	if key == "" {
		return Result{Key: key, Err: fmt.Errorf("key must not be empty")}
	}
	if !validKeyPattern.MatchString(key) {
		return Result{Key: key, Err: fmt.Errorf("key %q contains invalid characters", key)}
	}
	var warn string
	if v.warnOnLower && strings.ToUpper(key) != key {
		warn = fmt.Sprintf("key %q contains lowercase letters", key)
	}
	return Result{Key: key, Warning: warn}
}

// ValidateValue checks whether a value satisfies length constraints.
func (v *Validator) ValidateValue(key, value string) Result {
	if len(value) > v.maxValueLen {
		return Result{
			Key: key,
			Err: fmt.Errorf("value for key %q exceeds maximum length of %d", key, v.maxValueLen),
		}
	}
	return Result{Key: key}
}

// ValidateAll validates all key/value pairs and returns any results with
// errors or warnings. An empty slice means everything passed.
func (v *Validator) ValidateAll(secrets map[string]string) []Result {
	var results []Result
	for k, val := range secrets {
		if r := v.ValidateKey(k); !r.IsValid() || r.Warning != "" {
			results = append(results, r)
		}
		if r := v.ValidateValue(k, val); !r.IsValid() {
			results = append(results, r)
		}
	}
	return results
}
