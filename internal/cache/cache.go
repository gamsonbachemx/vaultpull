package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a cached set of secrets with metadata.
type Entry struct {
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
	TTL       time.Duration     `json:"ttl"`
}

// IsExpired returns true if the cache entry has exceeded its TTL.
func (e *Entry) IsExpired() bool {
	if e.TTL <= 0 {
		return true
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// Cache manages on-disk caching of Vault secrets.
type Cache struct {
	dir string
	ttl time.Duration
}

// New creates a new Cache that stores entries in dir with the given TTL.
func New(dir string, ttl time.Duration) (*Cache, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("cache: create directory: %w", err)
	}
	return &Cache{dir: dir, ttl: ttl}, nil
}

// key derives a filesystem-safe cache key from a Vault path.
func (c *Cache) key(path string) string {
	h := sha256.Sum256([]byte(path))
	return hex.EncodeToString(h[:])
}

// Get retrieves a non-expired cache entry for the given path.
// Returns nil, nil when no valid entry exists.
func (c *Cache) Get(path string) (*Entry, error) {
	file := filepath.Join(c.dir, c.key(path)+".json")
	data, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache: read: %w", err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("cache: unmarshal: %w", err)
	}
	if entry.IsExpired() {
		_ = os.Remove(file)
		return nil, nil
	}
	return &entry, nil
}

// Set writes secrets to the cache for the given path.
func (c *Cache) Set(path string, secrets map[string]string) error {
	entry := Entry{
		Secrets:   secrets,
		FetchedAt: time.Now(),
		TTL:       c.ttl,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	file := filepath.Join(c.dir, c.key(path)+".json")
	if err := os.WriteFile(file, data, 0600); err != nil {
		return fmt.Errorf("cache: write: %w", err)
	}
	return nil
}

// Invalidate removes the cache entry for the given path.
func (c *Cache) Invalidate(path string) error {
	file := filepath.Join(c.dir, c.key(path)+".json")
	err := os.Remove(file)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
