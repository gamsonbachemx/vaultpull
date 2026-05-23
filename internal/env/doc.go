// Package env provides utilities for reading, writing, parsing, merging,
// sanitizing, and exporting environment variable files (.env).
//
// # Reader
//
// NewReader reads an existing .env file from disk into a map.
//
// # Writer
//
// NewWriter writes a map of secrets to a .env file, formatting keys via
// FormatKey and quoting values that contain special characters.
//
// # Parser
//
// NewParser parses .env-formatted content from an io.Reader, handling
// comments, blank lines, quoted values, and the optional `export` prefix.
//
// # Merger
//
// Merge combines an incoming map of secrets with an existing map using
// either an overwrite or keep-existing strategy.
//
// # Sanitizer
//
// NewSanitizer validates and normalises keys and values, optionally
// uppercasing keys and stripping control characters from values.
//
// # Exporter
//
// NewExporter writes secrets as shell `export KEY=VALUE` statements
// suitable for evaluation in a shell session:
//
//	eval $(vaultpull export)
//
// Values that contain shell-special characters are wrapped in single
// quotes; embedded single quotes are escaped using the '\'' idiom.
package env
