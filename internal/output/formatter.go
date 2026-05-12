package output

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Format represents the output format for secrets.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Formatter writes secrets in a specific format.
type Formatter struct {
	format Format
	out    io.Writer
}

// New creates a new Formatter with the given format.
// If out is nil, os.Stdout is used.
func New(format Format, out io.Writer) *Formatter {
	if out == nil {
		out = os.Stdout
	}
	return &Formatter{format: format, out: out}
}

// Write outputs secrets in the configured format.
func (f *Formatter) Write(secrets map[string]string) error {
	switch f.format {
	case FormatExport:
		return f.writeExport(secrets)
	case FormatJSON:
		return f.writeJSON(secrets)
	default:
		return f.writeDotenv(secrets)
	}
}

func (f *Formatter) writeDotenv(secrets map[string]string) error {
	for k, v := range secrets {
		if _, err := fmt.Fprintf(f.out, "%s=%s\n", k, quoteIfNeeded(v)); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) writeExport(secrets map[string]string) error {
	for k, v := range secrets {
		if _, err := fmt.Fprintf(f.out, "export %s=%s\n", k, quoteIfNeeded(v)); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) writeJSON(secrets map[string]string) error {
	pairs := make([]string, 0, len(secrets))
	for k, v := range secrets {
		pairs = append(pairs, fmt.Sprintf(`  %q: %q`, k, v))
	}
	_, err := fmt.Fprintf(f.out, "{\n%s\n}\n", strings.Join(pairs, ",\n"))
	return err
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n#$") {
		return fmt.Sprintf(`"%s"`, strings.ReplaceAll(v, `"`, `\"`))
	}
	return v
}
