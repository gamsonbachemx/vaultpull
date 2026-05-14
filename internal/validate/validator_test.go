package validate_test

import (
	"strings"
	"testing"

	"vaultpull/internal/validate"
)

func TestValidateKey_Valid(t *testing.T) {
	v := validate.New()
	for _, key := range []string{"FOO", "_BAR", "A1_B2", "__PRIVATE"} {
		r := v.ValidateKey(key)
		if !r.IsValid() {
			t.Errorf("expected %q to be valid, got err: %v", key, r.Err)
		}
	}
}

func TestValidateKey_Invalid(t *testing.T) {
	v := validate.New()
	for _, key := range []string{"", "1START", "HAS-DASH", "HAS SPACE", "HAS.DOT"} {
		r := v.ValidateKey(key)
		if r.IsValid() {
			t.Errorf("expected %q to be invalid", key)
		}
	}
}

func TestValidateKey_WarnOnLowercase(t *testing.T) {
	v := validate.New(validate.WithWarnOnLowercase())

	r := v.ValidateKey("myKey")
	if !r.IsValid() {
		t.Fatalf("expected valid result, got err: %v", r.Err)
	}
	if r.Warning == "" {
		t.Error("expected a warning for lowercase key, got none")
	}

	r2 := v.ValidateKey("MY_KEY")
	if r2.Warning != "" {
		t.Errorf("expected no warning for uppercase key, got: %s", r2.Warning)
	}
}

func TestValidateValue_WithinLimit(t *testing.T) {
	v := validate.New(validate.WithMaxValueLen(10))
	r := v.ValidateValue("KEY", "short")
	if !r.IsValid() {
		t.Errorf("expected valid result, got: %v", r.Err)
	}
}

func TestValidateValue_ExceedsLimit(t *testing.T) {
	v := validate.New(validate.WithMaxValueLen(5))
	r := v.ValidateValue("KEY", "toolongvalue")
	if r.IsValid() {
		t.Error("expected invalid result for oversized value")
	}
	if !strings.Contains(r.Err.Error(), "exceeds maximum length") {
		t.Errorf("unexpected error message: %v", r.Err)
	}
}

func TestValidateAll_NoIssues(t *testing.T) {
	v := validate.New()
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	results := v.ValidateAll(secrets)
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestValidateAll_WithErrors(t *testing.T) {
	v := validate.New(validate.WithMaxValueLen(4))
	secrets := map[string]string{
		"GOOD_KEY": "ok",
		"BAD-KEY":  "val",
		"LONG_VAL": "toolongvalue",
	}
	results := v.ValidateAll(secrets)
	if len(results) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(results))
	}
}
