package resolve_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/resolve"
)

func TestBuild_KVv2(t *testing.T) {
	b := resolve.NewPathBuilder("secret", 2)
	got := b.Build("myapp/config")
	if want := "secret/data/myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuild_KVv1(t *testing.T) {
	b := resolve.NewPathBuilder("secret", 1)
	got := b.Build("myapp/config")
	if want := "secret/myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuild_LeadingSlash(t *testing.T) {
	b := resolve.NewPathBuilder("/secret/", 2)
	got := b.Build("/myapp/config/")
	if want := "secret/data/myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildMetadata_KVv2(t *testing.T) {
	b := resolve.NewPathBuilder("secret", 2)
	got := b.BuildMetadata("myapp/config")
	if want := "secret/metadata/myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildMetadata_KVv1(t *testing.T) {
	b := resolve.NewPathBuilder("secret", 1)
	got := b.BuildMetadata("myapp/config")
	if want := "secret/myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripMount_KVv2(t *testing.T) {
	b := resolve.NewPathBuilder("secret", 2)
	got := b.StripMount("secret/data/myapp/config")
	if want := "myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripMount_KVv1(t *testing.T) {
	b := resolve.NewPathBuilder("secret", 1)
	got := b.StripMount("secret/myapp/config")
	if want := "myapp/config"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNewPathBuilder_InvalidVersion(t *testing.T) {
	// Invalid version should default to v2.
	b := resolve.NewPathBuilder("secret", 99)
	got := b.Build("myapp")
	if want := "secret/data/myapp"; got != want {
		t.Errorf("expected v2 default, got %q", got)
	}
}
