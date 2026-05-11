package filter

import "strings"

// Filter holds namespace filtering configuration.
type Filter struct {
	namespaces []string
}

// New creates a new Filter with the given namespace prefixes.
// An empty slice means no filtering (all keys pass).
func New(namespaces []string) *Filter {
	return &Filter{namespaces: namespaces}
}

// Match reports whether the given key matches any of the configured
// namespace prefixes. If no namespaces are configured, all keys match.
func (f *Filter) Match(key string) bool {
	if len(f.namespaces) == 0 {
		return true
	}
	for _, ns := range f.namespaces {
		if strings.HasPrefix(key, ns) {
			return true
		}
	}
	return false
}

// Apply filters a map of secrets, returning only entries whose keys
// match one of the configured namespace prefixes.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	if len(f.namespaces) == 0 {
		return secrets
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if f.Match(k) {
			result[k] = v
		}
	}
	return result
}

// StripPrefix removes the matched namespace prefix from a key.
// If no prefix matches, the original key is returned unchanged.
func (f *Filter) StripPrefix(key string) string {
	for _, ns := range f.namespaces {
		if strings.HasPrefix(key, ns) {
			return strings.TrimPrefix(key, ns)
		}
	}
	return key
}
