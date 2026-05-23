package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Exporter writes secrets as shell export statements suitable for
// sourcing directly into a shell session (e.g. eval $(vaultpull export)).
type Exporter struct {
	w io.Writer
}

// NewExporter creates an Exporter that writes to w.
func NewExporter(w io.Writer) *Exporter {
	return &Exporter{w: w}
}

// Export writes each key/value pair as a shell export statement.
// Keys are sorted for deterministic output.
func (e *Exporter) Export(secrets map[string]string) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := secrets[k]
		formatted := FormatKey(k)
		escaped := escapeShellValue(v)
		if _, err := fmt.Fprintf(e.w, "export %s=%s\n", formatted, escaped); err != nil {
			return fmt.Errorf("exporter: write key %q: %w", formatted, err)
		}
	}
	return nil
}

// escapeShellValue wraps a value in single quotes, escaping any embedded
// single quotes using the '\'' idiom.
func escapeShellValue(v string) string {
	if !needsQuoting(v) {
		return v
	}
	escaped := strings.ReplaceAll(v, "'", "'\\''")
	return "'" + escaped + "'"
}

// needsQuoting reports whether v contains characters that require quoting
// when used in a shell assignment.
func needsQuoting(v string) bool {
	for _, c := range v {
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case c == '_' || c == '-' || c == '.' || c == '/' || c == ':':
		default:
			return true
		}
	}
	return false
}
