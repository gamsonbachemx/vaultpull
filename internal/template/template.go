// Package template provides .env file template rendering with variable substitution.
package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Renderer renders .env templates using secrets from Vault.
type Renderer struct {
	delimLeft  string
	delimRight string
}

// Option configures a Renderer.
type Option func(*Renderer)

// WithDelimiters sets custom template delimiters.
func WithDelimiters(left, right string) Option {
	return func(r *Renderer) {
		r.delimLeft = left
		r.delimRight = right
	}
}

// New creates a new Renderer with optional configuration.
func New(opts ...Option) *Renderer {
	r := &Renderer{
		delimLeft:  "{{",
		delimRight: "}}",
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Render applies secrets to a template string and returns the rendered output.
// Secrets are accessible via the .Secrets map in the template.
func (r *Renderer) Render(tmplStr string, secrets map[string]string) (string, error) {
	if tmplStr == "" {
		return "", nil
	}

	funcMap := template.FuncMap{
		"required": func(key string, secrets map[string]string) (string, error) {
			v, ok := secrets[key]
			if !ok || v == "" {
				return "", fmt.Errorf("required secret %q is missing or empty", key)
			}
			return v, nil
		},
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	t, err := template.New("env").
		Delims(r.delimLeft, r.delimRight).
		Funcs(funcMap).
		Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]interface{}{"Secrets": secrets}); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}

	return buf.String(), nil
}
