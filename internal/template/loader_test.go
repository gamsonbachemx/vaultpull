package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/template"
)

func writeTemplateFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	return path
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	writeTemplateFile(t, dir, "app.env.tmpl", `DB={{ index .Secrets "DB" }}`)

	l := template.NewLoader(dir)
	content, err := l.Load("app.env.tmpl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content == "" {
		t.Error("expected non-empty content")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	l := template.NewLoader(t.TempDir())
	_, err := l.Load("missing.tmpl")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_AbsolutePath(t *testing.T) {
	dir := t.TempDir()
	path := writeTemplateFile(t, dir, "abs.env.tmpl", "KEY=value")

	l := template.NewLoader("/some/other/dir")
	content, err := l.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != "KEY=value" {
		t.Errorf("got %q", content)
	}
}

func TestLoad_RenderIntegration(t *testing.T) {
	dir := t.TempDir()
	writeTemplateFile(t, dir, "svc.env.tmpl",
		`PORT={{ index .Secrets "PORT" }}
HOST={{ default "localhost" (index .Secrets "HOST") }}`)

	l := template.NewLoader(dir)
	content, err := l.Load("svc.env.tmpl")
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	r := template.New()
	out, err := r.Render(content, map[string]string{"PORT": "8080"})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if out != "PORT=8080\nHOST=localhost" {
		t.Errorf("got %q", out)
	}
}
