// Package backup provides functionality to create backups of existing .env files
// before overwriting them during a sync operation.
package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Manager handles backup creation for env files.
type Manager struct {
	backupDir string
	clock     func() time.Time
}

// New creates a new backup Manager. If backupDir is empty, backups are
// placed alongside the original file with a timestamp suffix.
func New(backupDir string) *Manager {
	return &Manager{
		backupDir: backupDir,
		clock:     time.Now,
	}
}

// Backup copies src to a timestamped backup file. It returns the path of
// the created backup, or an empty string if src does not exist.
func (m *Manager) Backup(src string) (string, error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return "", nil
	}

	timestamp := m.clock().UTC().Format("20060102T150405Z")
	base := filepath.Base(src)
	backupName := fmt.Sprintf("%s.%s.bak", base, timestamp)

	var dest string
	if m.backupDir != "" {
		if err := os.MkdirAll(m.backupDir, 0o700); err != nil {
			return "", fmt.Errorf("backup: create directory: %w", err)
		}
		dest = filepath.Join(m.backupDir, backupName)
	} else {
		dest = filepath.Join(filepath.Dir(src), backupName)
	}

	if err := copyFile(src, dest); err != nil {
		return "", fmt.Errorf("backup: copy %s -> %s: %w", src, dest, err)
	}

	return dest, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
