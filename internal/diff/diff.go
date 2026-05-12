package diff

import "sort"

// Result holds the comparison between existing and new secrets.
type Result struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]string
	Unchanged map[string]string
}

// HasChanges returns true if there are any additions, removals, or changes.
func (r *Result) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// Keys returns a sorted list of all keys across all categories.
func (r *Result) Keys() []string {
	seen := make(map[string]struct{})
	for k := range r.Added {
		seen[k] = struct{}{}
	}
	for k := range r.Removed {
		seen[k] = struct{}{}
	}
	for k := range r.Changed {
		seen[k] = struct{}{}
	}
	for k := range r.Unchanged {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Compare computes the diff between existing (old) and incoming (new) secret maps.
func Compare(existing, incoming map[string]string) *Result {
	r := &Result{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, newVal := range incoming {
		oldVal, exists := existing[k]
		if !exists {
			r.Added[k] = newVal
		} else if oldVal != newVal {
			r.Changed[k] = newVal
		} else {
			r.Unchanged[k] = newVal
		}
	}

	for k, oldVal := range existing {
		if _, exists := incoming[k]; !exists {
			r.Removed[k] = oldVal
		}
	}

	return r
}
