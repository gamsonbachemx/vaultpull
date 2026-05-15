package env

import (
	"testing"
)

func TestMerge_NilExisting(t *testing.T) {
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res := Merge(nil, incoming, StrategyOverwrite)
	if len(res.Merged) != 2 {
		t.Errorf("expected 2 merged keys, got %d", len(res.Merged))
	}
	if len(res.Added) != 2 {
		t.Errorf("expected 2 added keys, got %d", len(res.Added))
	}
	if len(res.Updated) != 0 || len(res.Skipped) != 0 {
		t.Errorf("expected no updates or skips")
	}
}

func TestMerge_OverwriteStrategy(t *testing.T) {
	existing := map[string]string{"FOO": "old", "KEEP": "same"}
	incoming := map[string]string{"FOO": "new", "KEEP": "same", "NEW": "val"}
	res := Merge(existing, incoming, StrategyOverwrite)

	if res.Merged["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %q", res.Merged["FOO"])
	}
	if res.Merged["KEEP"] != "same" {
		t.Errorf("expected KEEP=same, got %q", res.Merged["KEEP"])
	}
	if res.Merged["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %q", res.Merged["NEW"])
	}
	if len(res.Updated) != 1 || res.Updated[0] != "FOO" {
		t.Errorf("expected [FOO] updated, got %v", res.Updated)
	}
	if len(res.Added) != 1 || res.Added[0] != "NEW" {
		t.Errorf("expected [NEW] added, got %v", res.Added)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "KEEP" {
		t.Errorf("expected [KEEP] skipped, got %v", res.Skipped)
	}
}

func TestMerge_KeepExistingStrategy(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new", "BAR": "baz"}
	res := Merge(existing, incoming, StrategyKeepExisting)

	if res.Merged["FOO"] != "old" {
		t.Errorf("expected FOO=old (preserved), got %q", res.Merged["FOO"])
	}
	if res.Merged["BAR"] != "baz" {
		t.Errorf("expected BAR=baz, got %q", res.Merged["BAR"])
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped key, got %d", len(res.Skipped))
	}
	if len(res.Added) != 1 {
		t.Errorf("expected 1 added key, got %d", len(res.Added))
	}
}

func TestMerge_EmptyIncoming(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	res := Merge(existing, map[string]string{}, StrategyOverwrite)
	if len(res.Merged) != 1 {
		t.Errorf("expected 1 merged key, got %d", len(res.Merged))
	}
	if len(res.Added) != 0 || len(res.Updated) != 0 || len(res.Skipped) != 0 {
		t.Errorf("expected no changes")
	}
}
