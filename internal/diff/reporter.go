package diff

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Reporter writes a human-readable diff summary to a writer.
type Reporter struct {
	w io.Writer
}

// NewReporter creates a Reporter that writes to w.
// If w is nil, os.Stdout is used.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Report writes a summary of the diff result.
func (r *Reporter) Report(result *Result) {
	if !result.HasChanges() {
		fmt.Fprintln(r.w, "No changes detected.")
		return
	}

	if len(result.Added) > 0 {
		fmt.Fprintln(r.w, "Added:")
		for _, k := range sortedKeys(result.Added) {
			fmt.Fprintf(r.w, "  + %s\n", k)
		}
	}

	if len(result.Removed) > 0 {
		fmt.Fprintln(r.w, "Removed:")
		for _, k := range sortedKeys(result.Removed) {
			fmt.Fprintf(r.w, "  - %s\n", k)
		}
	}

	if len(result.Changed) > 0 {
		fmt.Fprintln(r.w, "Changed:")
		for _, k := range sortedKeys(result.Changed) {
			fmt.Fprintf(r.w, "  ~ %s\n", k)
		}
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
