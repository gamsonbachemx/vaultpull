package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds all runtime configuration for vaultpull.
type Config struct {
	VaultAddress string
	VaultToken   string
	SecretPath   string
	OutputFile   string
	Namespaces   []string
	StripPrefix  bool
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	addr := getEnv("VAULT_ADDR", "")
	if addr == "" {
		return nil, errors.New("VAULT_ADDR is required")
	}
	token := getEnv("VAULT_TOKEN", "")
	if token == "" {
		return nil, errors.New("VAULT_TOKEN is required")
	}

	var namespaces []string
	if ns := getEnv("VAULT_NAMESPACES", ""); ns != "" {
		for _, n := range strings.Split(ns, ",") {
			n = strings.TrimSpace(n)
			if n != "" {
				namespaces = append(namespaces, n)
			}
		}
	}

	return &Config{
		VaultAddress: addr,
		VaultToken:   token,
		SecretPath:   getEnv("VAULT_SECRET_PATH", "secret/data/app"),
		OutputFile:   getEnv("OUTPUT_FILE", ".env"),
		Namespaces:   namespaces,
		StripPrefix:  getEnv("STRIP_PREFIX", "false") == "true",
	}, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
