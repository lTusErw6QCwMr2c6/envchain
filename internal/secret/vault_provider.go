package secret

import (
	"fmt"
	"os"
	"strings"
)

// VaultProvider implements Provider using HashiCorp Vault via environment
// variables for token/address configuration. This is a lightweight client
// that reads/writes secrets to Vault's KV v2 secrets engine.
type VaultProvider struct {
	address   string
	token     string
	mountPath string
}

// NewVaultProvider creates a VaultProvider using VAULT_ADDR and VAULT_TOKEN
// environment variables. mountPath is the KV v2 mount (e.g. "secret").
func NewVaultProvider(mountPath string) (*VaultProvider, error) {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("vault: VAULT_ADDR environment variable not set")
	}
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("vault: VAULT_TOKEN environment variable not set")
	}
	if mountPath == "" {
		mountPath = "secret"
	}
	return &VaultProvider{
		address:   strings.TrimRight(addr, "/"),
		token:     token,
		mountPath: mountPath,
	}, nil
}

// secretPath returns the Vault KV v2 data path for a given key.
func (v *VaultProvider) secretPath(key string) string {
	return fmt.Sprintf("%s/v1/%s/data/envchain/%s", v.address, v.mountPath, key)
}

// Set stores a secret value in Vault under the envchain namespace.
func (v *VaultProvider) Set(key, value string) error {
	body := fmt.Sprintf(`{"data":{"value":%q}}`, value)
	req, err := newHTTPRequest("POST", v.secretPath(key), v.token, body)
	if err != nil {
		return fmt.Errorf("vault: set %q: %w", key, err)
	}
	if err := doHTTPRequest(req, 200); err != nil {
		return fmt.Errorf("vault: set %q: %w", key, err)
	}
	return nil
}

// Get retrieves a secret value from Vault.
func (v *VaultProvider) Get(key string) (string, error) {
	req, err := newHTTPRequest("GET", v.secretPath(key), v.token, "")
	if err != nil {
		return "", fmt.Errorf("vault: get %q: %w", key, err)
	}
	val, err := doHTTPRequestValue(req)
	if err != nil {
		return "", fmt.Errorf("vault: get %q: %w", key, err)
	}
	return val, nil
}

// Delete removes a secret from Vault by posting to the metadata endpoint.
func (v *VaultProvider) Delete(key string) error {
	metaPath := fmt.Sprintf("%s/v1/%s/metadata/envchain/%s", v.address, v.mountPath, key)
	req, err := newHTTPRequest("DELETE", metaPath, v.token, "")
	if err != nil {
		return fmt.Errorf("vault: delete %q: %w", key, err)
	}
	if err := doHTTPRequest(req, 204); err != nil {
		return fmt.Errorf("vault: delete %q: %w", key, err)
	}
	return nil
}
