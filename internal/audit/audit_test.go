package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"vaultpull/internal/audit"
)

func TestNew_DefaultsToStderr(t *testing.T) {
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLog_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Log(audit.EventSync, "test message", map[string]any{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if e.Type != audit.EventSync {
		t.Errorf("expected type %q, got %q", audit.EventSync, e.Type)
	}
	if e.Message != "test message" {
		t.Errorf("expected message %q, got %q", "test message", e.Message)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLog_NoDetails(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Log(audit.EventRead, "read event", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "read event") {
		t.Error("expected output to contain message")
	}
}

func TestSync_WritesDetails(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Sync("/secret/app", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Details["path"] != "/secret/app" {
		t.Errorf("expected path /secret/app, got %v", e.Details["path"])
	}
	if int(e.Details["count"].(float64)) != 5 {
		t.Errorf("expected count 5, got %v", e.Details["count"])
	}
}

func TestRead_WritesPath(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Read("/secret/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "/secret/db") {
		t.Error("expected output to contain vault path")
	}
}
