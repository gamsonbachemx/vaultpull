// Package resolve provides secret path resolution with support for
// variable interpolation and path aliasing.
package resolve

import (
	"fmt"
	"regexp"
	"strings"
)

var placeholderRe = regexp.MustCompile(`\$\{([^}]+)\}`)

// Resolver resolves Vault secret paths, supporting variable interpolation
// using ${VAR} syntax and static path aliases.
type Resolver struct {
	vars    map[string]string
	aliases map[string]string
}

// Option configures a Resolver.
type Option func(*Resolver)

// WithVars sets variables available for interpolation.
func WithVars(vars map[string]string) Option {
	return func(r *Resolver) {
		for k, v := range vars {
			r.vars[k] = v
		}
	}
}

// WithAliases registers path aliases (alias -> canonical path).
func WithAliases(aliases map[string]string) Option {
	return func(r *Resolver) {
		for k, v := range aliases {
			r.aliases[k] = v
		}
	}
}

// New creates a new Resolver with the given options.
func New(opts ...Option) *Resolver {
	r := &Resolver{
		vars:    make(map[string]string),
		aliases: make(map[string]string),
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Resolve takes a raw path (possibly containing aliases and ${VAR} placeholders)
// and returns the fully resolved Vault path.
func (r *Resolver) Resolve(path string) (string, error) {
	// Expand aliases first (exact match on full path).
	if canonical, ok := r.aliases[path]; ok {
		path = canonical
	}

	// Interpolate ${VAR} placeholders.
	var missing []string
	resolved := placeholderRe.ReplaceAllStringFunc(path, func(match string) string {
		key := match[2 : len(match)-1] // strip ${ and }
		if val, ok := r.vars[key]; ok {
			return val
		}
		missing = append(missing, key)
		return match
	})

	if len(missing) > 0 {
		return "", fmt.Errorf("resolve: undefined variables: %s", strings.Join(missing, ", "))
	}

	return resolved, nil
}

// ResolveAll resolves a slice of paths, returning all errors encountered.
func (r *Resolver) ResolveAll(paths []string) ([]string, error) {
	out := make([]string, 0, len(paths))
	var errs []string
	for _, p := range paths {
		resolved, err := r.Resolve(p)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		out = append(out, resolved)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return out, nil
}
