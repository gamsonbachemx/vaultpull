// Package hook provides pre/post sync lifecycle hooks for vaultpull.
// Hooks allow users to run shell commands before or after a sync operation.
package hook

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Hook represents a shell command to run at a lifecycle point.
type Hook struct {
	Command string
	Args    []string
}

// Runner executes lifecycle hooks.
type Runner struct {
	pre  []Hook
	post []Hook
}

// Config holds pre and post hook command strings.
type Config struct {
	Pre  []string
	Post []string
}

// New creates a Runner from a Config.
// Each command string is split on whitespace into command + args.
func New(cfg Config) (*Runner, error) {
	pre, err := parseHooks(cfg.Pre)
	if err != nil {
		return nil, fmt.Errorf("pre-hook: %w", err)
	}
	post, err := parseHooks(cfg.Post)
	if err != nil {
		return nil, fmt.Errorf("post-hook: %w", err)
	}
	return &Runner{pre: pre, post: post}, nil
}

// RunPre executes all pre-sync hooks in order.
func (r *Runner) RunPre(ctx context.Context) error {
	return r.run(ctx, r.pre)
}

// RunPost executes all post-sync hooks in order.
func (r *Runner) RunPost(ctx context.Context) error {
	return r.run(ctx, r.post)
}

func (r *Runner) run(ctx context.Context, hooks []Hook) error {
	for _, h := range hooks {
		cmd := exec.CommandContext(ctx, h.Command, h.Args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("hook %q failed: %w\noutput: %s", h.Command, err, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func parseHooks(cmds []string) ([]Hook, error) {
	hooks := make([]Hook, 0, len(cmds))
	for _, c := range cmds {
		parts := strings.Fields(c)
		if len(parts) == 0 {
			return nil, fmt.Errorf("empty command string")
		}
		hooks = append(hooks, Hook{Command: parts[0], Args: parts[1:]})
	}
	return hooks, nil
}
