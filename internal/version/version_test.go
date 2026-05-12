package version_test

import (
	"strings"
	"testing"

	"vaultpull/internal/version"
)

func TestGet_ReturnsInfo(t *testing.T) {
	info := version.Get()

	if info.Version == "" {
		t.Error("expected non-empty Version")
	}
	if info.Commit == "" {
		t.Error("expected non-empty Commit")
	}
	if info.BuildDate == "" {
		t.Error("expected non-empty BuildDate")
	}
}

func TestGet_DefaultValues(t *testing.T) {
	info := version.Get()

	if info.Version != "dev" {
		t.Errorf("expected default version 'dev', got %q", info.Version)
	}
	if info.Commit != "none" {
		t.Errorf("expected default commit 'none', got %q", info.Commit)
	}
	if info.BuildDate != "unknown" {
		t.Errorf("expected default build_date 'unknown', got %q", info.BuildDate)
	}
}

func TestInfo_String(t *testing.T) {
	info := version.Info{
		Version:   "1.2.3",
		Commit:    "abc1234",
		BuildDate: "2024-01-15",
	}

	s := info.String()

	if !strings.Contains(s, "vaultpull") {
		t.Errorf("expected string to contain 'vaultpull', got %q", s)
	}
	if !strings.Contains(s, "1.2.3") {
		t.Errorf("expected string to contain version '1.2.3', got %q", s)
	}
	if !strings.Contains(s, "abc1234") {
		t.Errorf("expected string to contain commit 'abc1234', got %q", s)
	}
	if !strings.Contains(s, "2024-01-15") {
		t.Errorf("expected string to contain build date '2024-01-15', got %q", s)
	}
}

func TestInfo_Short(t *testing.T) {
	info := version.Info{
		Version:   "1.2.3",
		Commit:    "abc1234",
		BuildDate: "2024-01-15",
	}

	short := info.Short()

	if short != "1.2.3" {
		t.Errorf("expected short version '1.2.3', got %q", short)
	}
}
