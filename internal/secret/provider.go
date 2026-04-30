package secret

import "context"

// Provider defines the interface for secret backends.
// Implementations include keyring, env, vault, and AWS Secrets Manager.
type Provider interface {
	// Set stores a secret value for the given profile and key.
	Set(ctx context.Context, profile, key, value string) error

	// Get retrieves a secret value for the given profile and key.
	// Returns an error if the secret does not exist.
	Get(ctx context.Context, profile, key string) (string, error)

	// Delete removes the secret for the given profile and key.
	Delete(ctx context.Context, profile, key string) error
}

// ProviderType enumerates supported secret backends.
type ProviderType string

const (
	ProviderEnv     ProviderType = "env"
	ProviderKeyring ProviderType = "keyring"
	ProviderVault   ProviderType = "vault"
	ProviderAWS     ProviderType = "aws"
)

// ErrNotFound is returned when a secret does not exist in the backend.
type ErrNotFound struct {
	Profile string
	Key     string
}

func (e *ErrNotFound) Error() string {
	return "secret not found: profile=" + e.Profile + " key=" + e.Key
}
