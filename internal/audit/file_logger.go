package audit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileLogger wraps Logger to write audit events to a file.
type FileLogger struct {
	*Logger
	f *os.File
}

// NewFileLogger creates a Logger that appends audit events to the given file path.
// The directory is created if it does not exist.
func NewFileLogger(path string) (*FileLogger, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("audit: create directory %s: %w", dir, err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open file %s: %w", path, err)
	}

	return &FileLogger{
		Logger: New(f),
		f:      f,
	}, nil
}

// Close flushes and closes the underlying file.
func (fl *FileLogger) Close() error {
	if fl.f == nil {
		return nil
	}
	if err := fl.f.Sync(); err != nil {
		return fmt.Errorf("audit: sync file: %w", err)
	}
	return fl.f.Close()
}

// Discard returns a Logger that silently drops all events.
func Discard() *Logger {
	return New(io.Discard)
}
