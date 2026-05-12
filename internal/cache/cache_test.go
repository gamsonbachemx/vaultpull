package cache_test

import (
	"os"
	"testing"
	"time"

	"vaultpull/internal/cache"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "vaultpull-cache-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "vaultpull-newtest")
	defer os.RemoveAll(dir)
	_, err := cache.New(dir, time.Minute)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestGet_MissingEntry(t *testing.T) {
	c, _ := cache.New(tempDir(t), time.Minute)
	entry, err := c.Get("secret/data/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Fatal("expected nil entry for missing key")
	}
}

func TestSetAndGet_ValidEntry(t *testing.T) {
	c, _ := cache.New(tempDir(t), time.Minute)
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	if err := c.Set("secret/data/app", secrets); err != nil {
		t.Fatalf("Set: %v", err)
	}

	entry, err := c.Get("secret/data/app")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if entry == nil {
		t.Fatal("expected non-nil entry")
	}
	if entry.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", entry.Secrets["DB_HOST"])
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c, _ := cache.New(tempDir(t), time.Millisecond)
	_ = c.Set("secret/data/app", map[string]string{"KEY": "val"})
	time.Sleep(5 * time.Millisecond)

	entry, err := c.Get("secret/data/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Fatal("expected nil entry for expired cache")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c, _ := cache.New(tempDir(t), time.Minute)
	_ = c.Set("secret/data/app", map[string]string{"KEY": "val"})

	if err := c.Invalidate("secret/data/app"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}

	entry, _ := c.Get("secret/data/app")
	if entry != nil {
		t.Fatal("expected entry to be invalidated")
	}
}

func TestInvalidate_NonExistentKey(t *testing.T) {
	c, _ := cache.New(tempDir(t), time.Minute)
	if err := c.Invalidate("secret/data/missing"); err != nil {
		t.Fatalf("expected no error for missing key, got %v", err)
	}
}

func TestIsExpired_ZeroTTL(t *testing.T) {
	entry := &cache.Entry{FetchedAt: time.Now(), TTL: 0}
	if !entry.IsExpired() {
		t.Fatal("expected zero-TTL entry to be expired")
	}
}
