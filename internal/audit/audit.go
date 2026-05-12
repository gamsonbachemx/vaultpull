package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType represents the type of audit event.
type EventType string

const (
	EventSync    EventType = "sync"
	EventRead    EventType = "read"
	EventWrite   EventType = "write"
	EventFilter  EventType = "filter"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type"`
	Message   string    `json:"message"`
	Details   map[string]any `json:"details,omitempty"`
}

// Logger writes audit events to a destination.
type Logger struct {
	w io.Writer
}

// New creates a new audit Logger writing to w.
// If w is nil, os.Stderr is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{w: w}
}

// Log records an audit event.
func (l *Logger) Log(eventType EventType, message string, details map[string]any) error {
	e := Event{
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		Message:   message,
		Details:   details,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	if err != nil {
		return fmt.Errorf("audit: write event: %w", err)
	}
	return nil
}

// Sync logs a sync event.
func (l *Logger) Sync(path string, count int) error {
	return l.Log(EventSync, "secrets synced", map[string]any{
		"path":  path,
		"count": count,
	})
}

// Read logs a vault read event.
func (l *Logger) Read(path string) error {
	return l.Log(EventRead, "vault path read", map[string]any{
		"path": path,
	})
}
