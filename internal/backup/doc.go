// Package backup creates timestamped backups of existing .env files before
// vaultpull overwrites them during a sync. This allows users to recover
// previous secret values if needed.
//
// Basic usage:
//
//	m := backup.New(".vaultpull/backups") // or "" for same directory
//	dest, err := m.Backup(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if dest != "" {
//		fmt.Printf("backed up existing .env to %s\n", dest)
//	}
package backup

// NewWithClock is like New but accepts a custom clock function. Intended for
// testing purposes to produce deterministic backup filenames.
func NewWithClock(backupDir string, clock func() Time) *Manager {
	return &Manager{
		backupDir: backupDir,
		clock:     clock,
	}
}
