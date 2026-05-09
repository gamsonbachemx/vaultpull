package vault

import (
	"errors"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api       *vaultapi.Client
	namespace string
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(address, token, namespace string) (*Client, error) {
	if address == "" {
		return nil, errors.New("vault address must not be empty")
	}
	if token == "" {
		return nil, errors.New("vault token must not be empty")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault api client: %w", err)
	}
	api.SetToken(token)

	return &Client{
		api:       api,
		namespace: namespace,
	}, nil
}

// ReadSecrets reads all key-value pairs from the given KV v2 path,
// filtering keys by the configured namespace prefix.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read path %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at path %q", path)
	}

	// KV v2 wraps data under the "data" key.
	raw, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("unexpected secret format at path %q", path)
	}

	data, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("secret data at path %q is not a map", path)
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		if c.namespace != "" && !strings.HasPrefix(k, c.namespace) {
			continue
		}
		str, ok := v.(string)
		if !ok {
			continue
		}
		result[k] = str
	}

	return result, nil
}
