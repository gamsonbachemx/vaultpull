package backup_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"vaultpull/internal/backup"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func TestBackup_FileDoesNotExist(t *testing.T) {
	m := backup.New("")
	dest, err := m.Backup("/nonexistent/.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest != "" {
		t.Errorf("expected empty dest, got %q", dest)
	}
}

func TestBackup_SameDirectory(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	writeFile(t, src, "SECRET=hello\n")

	m := backup.New("")
	dest, err := m.Backup(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest == "" {
		t.Fatal("expected non-empty dest")
	}
	if !strings.HasSuffix(dest, ".bak") {
		t.Errorf("expected .bak suffix, got %q", dest)
	}
	if filepath.Dir(dest) != dir {
		t.Errorf("expected backup in %q, got %q", dir, filepath.Dir(dest))
	}
	data, _ := os.ReadFile(dest)
	if string(data) != "SECRET=hello\n" {
		t.Errorf("backup content mismatch: %q", data)
	}
}

func TestBackup_CustomDirectory(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	writeFile(t, src, "TOKEN=abc\n")

	backupDir := filepath.Join(dir, "backups")
	m := backup.New(backupDir)
	dest, err := m.Backup(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Dir(dest) != backupDir {
		t.Errorf("expected backup in %q, got %q", backupDir, filepath.Dir(dest))
	}
}

func TestBackup_TimestampInFilename(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	writeFile(t, src, "X=1\n")

	fixed := time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC)
	m := &backup.Manager{}
	// Use exported constructor and override clock via test helper approach:
	m2 := backup.NewWithClock("", func() time.Time { return fixed })
	dest, err := m2.Backup(src)
	_ = m
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(dest, "20240615T123045Z") {
		t.Errorf("expected timestamp in filename, got %q", dest)
	}
}
