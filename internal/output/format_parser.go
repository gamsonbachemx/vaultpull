package output

import (
	"fmt"
	"strings"
)

// ParseFormat converts a string to a Format type.
// Returns an error if the format string is not recognized.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "dotenv", "":
		return FormatDotenv, nil
	case "export":
		return FormatExport, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown output format %q: must be one of dotenv, export, json", s)
	}
}

// ValidFormats returns a slice of all supported format names.
func ValidFormats() []string {
	return []string{
		string(FormatDotenv),
		string(FormatExport),
		string(FormatJSON),
	}
}

// String implements the Stringer interface for Format.
func (f Format) String() string {
	return string(f)
}
