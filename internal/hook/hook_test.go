package hook_test

import (
	"context"
	"runtime"
	"testing"

	"github.com/user/vaultpull/internal/hook"
)

func TestNew_EmptyConfig(t *testing.T) {
	r, err := hook.New(hook.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestNew_InvalidPreHook(t *testing.T) {
	_, err := hook.New(hook.Config{Pre: []string{"", "echo hi"}})
	if err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestNew_InvalidPostHook(t *testing.T) {
	_, err := hook.New(hook.Config{Post: []string{""}}) 
	if err == nil {
		t.Fatal("expected error for empty post command")
	}
}

func TestRunPre_Success(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r, err := hook.New(hook.Config{
		Pre: []string{"echo pre-hook"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := r.RunPre(context.Background()); err != nil {
		t.Errorf("RunPre: %v", err)
	}
}

func TestRunPost_Success(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r, err := hook.New(hook.Config{
		Post: []string{"echo post-hook"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := r.RunPost(context.Background()); err != nil {
		t.Errorf("RunPost: %v", err)
	}
}

func TestRunPre_CommandFails(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r, err := hook.New(hook.Config{
		Pre: []string{"false"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := r.RunPre(context.Background()); err == nil {
		t.Error("expected error from failing command")
	}
}

func TestRunPre_MultipleHooks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r, err := hook.New(hook.Config{
		Pre: []string{"echo first", "echo second"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := r.RunPre(context.Background()); err != nil {
		t.Errorf("RunPre: %v", err)
	}
}

func TestRunPre_ContextCancelled(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r, err := hook.New(hook.Config{
		Pre: []string{"sleep 10"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := r.RunPre(ctx); err == nil {
		t.Error("expected error from cancelled context")
	}
}
