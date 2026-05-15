package template_test

import (
	"strings"
	"testing"

	"github.com/user/vaultpull/internal/template"
)

func TestRender_EmptyTemplate(t *testing.T) {
	r := template.New()
	out, err := r.Render("", map[string]string{"KEY": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}

func TestRender_BasicSubstitution(t *testing.T) {
	r := template.New()
	tmpl := `DB_HOST={{ index .Secrets "DB_HOST" }}`
	secrets := map[string]string{"DB_HOST": "localhost"}

	out, err := r.Render(tmpl, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "DB_HOST=localhost" {
		t.Errorf("got %q", out)
	}
}

func TestRender_RequiredMissing(t *testing.T) {
	r := template.New()
	tmpl := `{{ required "API_KEY" .Secrets }}`
	_, err := r.Render(tmpl, map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing required secret")
	}
	if !strings.Contains(err.Error(), "API_KEY") {
		t.Errorf("error should mention key name, got: %v", err)
	}
}

func TestRender_DefaultFunction(t *testing.T) {
	r := template.New()
	tmpl := `LOG_LEVEL={{ default "info" (index .Secrets "LOG_LEVEL") }}`
	out, err := r.Render(tmpl, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "LOG_LEVEL=info" {
		t.Errorf("got %q", out)
	}
}

func TestRender_UpperLowerFunctions(t *testing.T) {
	r := template.New()
	tmpl := `ENV={{ upper (index .Secrets "env") }}`
	out, err := r.Render(tmpl, map[string]string{"env": "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "ENV=PRODUCTION" {
		t.Errorf("got %q", out)
	}
}

func TestRender_CustomDelimiters(t *testing.T) {
	r := template.New(template.WithDelimiters("((", "))"))
	tmpl := `HOST=(( index .Secrets "HOST" ))`
	out, err := r.Render(tmpl, map[string]string{"HOST": "127.0.0.1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "HOST=127.0.0.1" {
		t.Errorf("got %q", out)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	r := template.New()
	_, err := r.Render(`{{ .Secrets["bad"] }`, map[string]string{})
	if err == nil {
		t.Fatal("expected parse error for invalid template")
	}
}
