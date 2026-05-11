package sync

import (
	"fmt"
	"log"

	"vaultpull/internal/config"
	"vaultpull/internal/env"
	"vaultpull/internal/vault"
)

// Syncer orchestrates reading secrets from Vault and writing them to a .env file.
type Syncer struct {
	cfg    *config.Config
	client *vault.Client
	writer *env.Writer
}

// New creates a new Syncer from the provided configuration.
func New(cfg *config.Config) (*Syncer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config must not be nil")
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	writer, err := env.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create env writer: %w", err)
	}

	return &Syncer{
		cfg:    cfg,
		client: client,
		writer: writer,
	}, nil
}

// Run performs the full sync: reads secrets from Vault and writes them to the output file.
func (s *Syncer) Run() error {
	log.Printf("reading secrets from path: %s", s.cfg.SecretPath)

	secrets, err := s.client.ReadSecrets(s.cfg.SecretPath, s.cfg.Namespace)
	if err != nil {
		return fmt.Errorf("failed to read secrets from vault: %w", err)
	}

	if len(secrets) == 0 {
		log.Println("no secrets found; output file will be empty")
	}

	log.Printf("writing %d secret(s) to %s", len(secrets), s.cfg.OutputFile)

	if err := s.writer.Write(secrets); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}

	log.Println("sync complete")
	return nil
}
