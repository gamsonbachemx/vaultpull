package diff_test

import (
	"testing"

	"github.com/vaultpull/internal/diff"
)

func TestCompare_EmptyBoth(t *testing.T) {
	r := diff.Compare(nil, nil)
	if r.HasChanges() {
		t.Error("expected no changes for empty maps")
	}
}

func TestCompare_AllAdded(t *testing.T) {
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := diff.Compare(nil, incoming)
	if len(r.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(r.Added))
	}
	if len(r.Removed) != 0 || len(r.Changed) != 0 {
		t.Error("expected no removed or changed")
	}
}

func TestCompare_AllRemoved(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := diff.Compare(existing, nil)
	if len(r.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(r.Removed))
	}
	if len(r.Added) != 0 || len(r.Changed) != 0 {
		t.Error("expected no added or changed")
	}
}

func TestCompare_ChangedAndUnchanged(t *testing.T) {
	existing := map[string]string{"FOO": "old", "KEEP": "same"}
	incoming := map[string]string{"FOO": "new", "KEEP": "same"}
	r := diff.Compare(existing, incoming)
	if len(r.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(r.Changed))
	}
	if r.Changed["FOO"] != "new" {
		t.Errorf("expected changed FOO=new, got %s", r.Changed["FOO"])
	}
	if len(r.Unchanged) != 1 {
		t.Errorf("expected 1 unchanged, got %d", len(r.Unchanged))
	}
	if !r.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_NoChanges(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	r := diff.Compare(secrets, secrets)
	if r.HasChanges() {
		t.Error("expected no changes when maps are identical")
	}
	if len(r.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(r.Unchanged))
	}
}

func TestResult_Keys_Sorted(t *testing.T) {
	existing := map[string]string{"Z": "1", "A": "2"}
	incoming := map[string]string{"Z": "changed", "M": "new"}
	r := diff.Compare(existing, incoming)
	keys := r.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("keys not sorted: %v", keys)
	}
}
