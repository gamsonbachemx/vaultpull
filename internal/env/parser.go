package env

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Parser reads and parses .env file content into a map of key-value pairs.
type Parser struct{}

// NewParser returns a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads from r and returns a map of key-value pairs.
// Lines beginning with '#' are treated as comments and skipped.
// Blank lines are skipped. Values may be optionally quoted.
func (p *Parser) Parse(r io.Reader) (map[string]string, error) {
	result := make(map[string]string)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip optional "export " prefix
		line = strings.TrimPrefix(line, "export ")

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("line %d: missing '=' in %q", lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		result[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return result, nil
}
