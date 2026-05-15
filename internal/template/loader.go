package template

import (
	"fmt"
	"os"
	"path/filepath"
)

// Loader reads template files from disk.
type Loader struct {
	baseDir string
}

// NewLoader creates a Loader rooted at baseDir.
func NewLoader(baseDir string) *Loader {
	return &Loader{baseDir: baseDir}
}

// Load reads a template file by name relative to the base directory.
// If name is an absolute path, baseDir is ignored.
func (l *Loader) Load(name string) (string, error) {
	path := name
	if !filepath.IsAbs(name) {
		path = filepath.Join(l.baseDir, name)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template file not found: %s", path)
		}
		return "", fmt.Errorf("read template file %s: %w", path, err)
	}

	return string(data), nil
}

// MustLoad reads a template file and panics on error.
func (l *Loader) MustLoad(name string) string {
	s, err := l.Load(name)
	if err != nil {
		panic(err)
	}
	return s
}
