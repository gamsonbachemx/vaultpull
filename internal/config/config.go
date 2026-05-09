package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds the runtime configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	Namespace  string
	OutputFile string
	Prefix     string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	addr := getEnv("VAULT_ADDR", "http://127.0.0.1:8200")
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, errors.New("VAULT_TOKEN environment variable is required")
	}

	namespace := getEnv("VAULTPULL_NAMESPACE", "secret")
	output := getEnv("VAULTPULL_OUTPUT", ".env")
	prefix := strings.ToUpper(getEnv("VAULTPULL_PREFIX", ""))

	return &Config{
		VaultAddr:  addr,
		VaultToken: token,
		Namespace:  namespace,
		OutputFile: output,
		Prefix:     prefix,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
