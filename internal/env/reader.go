package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Reader reads an existing .env file and returns its key-value pairs.
type Reader struct {
	path string
}

// NewReader creates a Reader for the given file path.
func NewReader(path string) *Reader {
	return &Reader{path: path}
}

// Read opens the .env file and returns a map of key-value pairs.
// If the file does not exist, an empty map is returned without error.
func (r *Reader) Read() (map[string]string, error) {
	f, err := os.Open(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("env reader: open %s: %w", r.path, err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip optional export prefix.
		line = strings.TrimPrefix(line, "export ")

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("env reader: %s line %d: missing '=' separator", r.path, lineNum)
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		// Strip surrounding quotes from value.
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		result[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env reader: scan %s: %w", r.path, err)
	}

	return result, nil
}
