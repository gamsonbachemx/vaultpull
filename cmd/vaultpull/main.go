// Package main is the entry point for the vaultpull CLI tool.
// It wires together configuration, vault access, filtering, syncing,
// auditing, caching, and output formatting into a single command.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/vaultpull/internal/audit"
	"github.com/user/vaultpull/internal/backup"
	"github.com/user/vaultpull/internal/cache"
	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/filter"
	"github.com/user/vaultpull/internal/hook"
	"github.com/user/vaultpull/internal/output"
	"github.com/user/vaultpull/internal/prompt"
	"github.com/user/vaultpull/internal/sync"
	"github.com/user/vaultpull/internal/vault"
	"github.com/user/vaultpull/internal/version"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// CLI flags
	var (
		showVersion = flag.Bool("version", false, "print version information and exit")
		envFile     = flag.String("out", ".env", "output .env file path")
		formatStr   = flag.String("format", "dotenv", "output format: dotenv, export, json")
		namespaces  = flag.String("namespaces", "", "comma-separated namespace prefixes to filter (e.g. APP,DB)")
		auditFile   = flag.String("audit-log", "", "path to write audit log (default: stderr)")
		cacheDir    = flag.String("cache-dir", "", "directory for caching secrets (empty disables cache)")
		backupDir   = flag.String("backup-dir", "", "directory to store .env backups (empty: same dir as out)")
		preHook     = flag.String("pre-hook", "", "command to run before syncing")
		postHook    = flag.String("post-hook", "", "command to run after syncing")
		yes         = flag.Bool("yes", false, "skip confirmation prompt")
		noCache     = flag.Bool("no-cache", false, "bypass cache and always fetch from Vault")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(version.Get())
		return nil
	}

	// Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// Override config with CLI flags where provided
	if *envFile != "" {
		cfg.EnvFile = *envFile
	}

	// Set up context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Set up audit logger
	var auditLogger *audit.Logger
	if *auditFile != "" {
		fl, err := audit.NewFileLogger(*auditFile)
		if err != nil {
			return fmt.Errorf("audit log: %w", err)
		}
		defer fl.Close()
		auditLogger = audit.New(fl)
	} else {
		auditLogger = audit.New(os.Stderr)
	}

	// Set up namespace filter
	ns := filter.New(*namespaces)

	// Set up output formatter
	fmt_, err := output.ParseFormat(*formatStr)
	if err != nil {
		return fmt.Errorf("format: %w", err)
	}
	outFile, err := os.OpenFile(cfg.EnvFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open output file %q: %w", cfg.EnvFile, err)
	}
	defer outFile.Close()
	formatter := output.New(outFile, fmt_)

	// Set up optional backup
	backer := backup.New(cfg.EnvFile, *backupDir)

	// Set up optional cache
	var secretCache *cache.Cache
	if *cacheDir != "" && !*noCache {
		secretCache, err = cache.New(*cacheDir)
		if err != nil {
			return fmt.Errorf("cache: %w", err)
		}
	}

	// Set up lifecycle hooks
	hooks, err := hook.New(*preHook, *postHook)
	if err != nil {
		return fmt.Errorf("hooks: %w", err)
	}

	// Confirm before overwriting existing file
	if !*yes {
		if _, statErr := os.Stat(cfg.EnvFile); statErr == nil {
			p := prompt.New(os.Stdin, os.Stdout)
			ok, promptErr := p.Confirm(fmt.Sprintf("Overwrite %q?", cfg.EnvFile), false)
			if promptErr != nil {
				return fmt.Errorf("prompt: %w", promptErr)
			}
			if !ok {
				fmt.Println("Aborted.")
				return nil
			}
		}
	}

	// Build Vault client
	vaultClient, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	// Build and run syncer
	syncer, err := sync.New(cfg, vaultClient, ns, formatter, auditLogger, backer, secretCache, hooks)
	if err != nil {
		return fmt.Errorf("syncer: %w", err)
	}

	if err := syncer.Run(ctx); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	fmt.Printf("Secrets written to %s\n", cfg.EnvFile)
	return nil
}
