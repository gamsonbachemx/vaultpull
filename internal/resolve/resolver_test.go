package resolve_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/resolve"
)

func TestResolve_NoInterpolation(t *testing.T) {
	r := resolve.New()
	got, err := r.Resolve("secret/myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/myapp/config" {
		t.Errorf("expected unchanged path, got %q", got)
	}
}

func TestResolve_WithVars(t *testing.T) {
	r := resolve.New(resolve.WithVars(map[string]string{
		"ENV": "production",
		"APP": "myapp",
	}))
	got, err := r.Resolve("secret/${ENV}/${APP}/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "secret/production/myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolve_MissingVar(t *testing.T) {
	r := resolve.New()
	_, err := r.Resolve("secret/${MISSING}/config")
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestResolve_WithAlias(t *testing.T) {
	r := resolve.New(resolve.WithAliases(map[string]string{
		"myapp-prod": "secret/production/myapp/config",
	}))
	got, err := r.Resolve("myapp-prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "secret/production/myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolve_AliasWithInterpolation(t *testing.T) {
	r := resolve.New(
		resolve.WithAliases(map[string]string{
			"myapp": "secret/${ENV}/myapp/config",
		}),
		resolve.WithVars(map[string]string{"ENV": "staging"}),
	)
	got, err := r.Resolve("myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "secret/staging/myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolveAll_Success(t *testing.T) {
	r := resolve.New(resolve.WithVars(map[string]string{"ENV": "dev"}))
	paths := []string{"secret/${ENV}/a", "secret/${ENV}/b"}
	got, err := r.ResolveAll(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0] != "secret/dev/a" || got[1] != "secret/dev/b" {
		t.Errorf("unexpected results: %v", got)
	}
}

func TestResolveAll_PartialError(t *testing.T) {
	r := resolve.New()
	paths := []string{"secret/valid", "secret/${MISSING}"}
	_, err := r.ResolveAll(paths)
	if err == nil {
		t.Fatal("expected error for missing variable in batch")
	}
}
