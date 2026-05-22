package resolve

import (
	"path"
	"strings"
)

// PathBuilder constructs Vault KV paths from component parts, handling
// both KV v1 and KV v2 mount conventions.
type PathBuilder struct {
	mount   string
	version int
}

// NewPathBuilder creates a PathBuilder for the given mount point.
// version should be 1 or 2 to match the KV secrets engine version.
func NewPathBuilder(mount string, version int) *PathBuilder {
	if version != 1 && version != 2 {
		version = 2
	}
	return &PathBuilder{
		mount:   strings.Trim(mount, "/"),
		version: version,
	}
}

// Build constructs the full API path for a secret at the given logical path.
// For KV v2, it inserts the "data/" infix required by the Vault HTTP API.
func (b *PathBuilder) Build(secretPath string) string {
	secretPath = strings.Trim(secretPath, "/")
	if b.version == 2 {
		return path.Join(b.mount, "data", secretPath)
	}
	return path.Join(b.mount, secretPath)
}

// BuildMetadata returns the metadata path for a KV v2 secret.
// Returns the plain path for KV v1 (no metadata concept).
func (b *PathBuilder) BuildMetadata(secretPath string) string {
	secretPath = strings.Trim(secretPath, "/")
	if b.version == 2 {
		return path.Join(b.mount, "metadata", secretPath)
	}
	return path.Join(b.mount, secretPath)
}

// StripMount removes the mount prefix from a full Vault path, returning
// the logical secret path. Useful when normalising paths from list responses.
func (b *PathBuilder) StripMount(fullPath string) string {
	prefix := b.mount + "/"
	if b.version == 2 {
		prefix = b.mount + "/data/"
	}
	return strings.TrimPrefix(fullPath, prefix)
}
