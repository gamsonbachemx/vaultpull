package output

import (
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
	}{
		{"dotenv", FormatDotenv},
		{"DOTENV", FormatDotenv},
		{"", FormatDotenv},
		{"export", FormatExport},
		{"Export", FormatExport},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("yaml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestValidFormats(t *testing.T) {
	formats := ValidFormats()
	if len(formats) != 3 {
		t.Errorf("expected 3 valid formats, got %d", len(formats))
	}
}

func TestFormat_String(t *testing.T) {
	if FormatDotenv.String() != "dotenv" {
		t.Errorf("expected dotenv, got %s", FormatDotenv.String())
	}
	if FormatExport.String() != "export" {
		t.Errorf("expected export, got %s", FormatExport.String())
	}
	if FormatJSON.String() != "json" {
		t.Errorf("expected json, got %s", FormatJSON.String())
	}
}
