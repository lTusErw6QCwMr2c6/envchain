package secret

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const keyringService = "envchain"

// KeyringProvider stores and retrieves secrets using the OS keyring.
type KeyringProvider struct {
	service string
}

// NewKeyringProvider creates a new KeyringProvider using the default service name.
func NewKeyringProvider() *KeyringProvider {
	return &KeyringProvider{service: keyringService}
}

// NewKeyringProviderWithService creates a KeyringProvider with a custom service name.
func NewKeyringProviderWithService(service string) *KeyringProvider {
	return &KeyringProvider{service: service}
}

// Set stores a secret in the OS keyring under the given profile and key.
func (k *KeyringProvider) Set(profile, key, value string) error {
	account := k.account(profile, key)
	if err := keyring.Set(k.service, account, value); err != nil {
		return fmt.Errorf("keyring set %q: %w", account, err)
	}
	return nil
}

// Get retrieves a secret from the OS keyring for the given profile and key.
func (k *KeyringProvider) Get(profile, key string) (string, error) {
	account := k.account(profile, key)
	value, err := keyring.Get(k.service, account)
	if err == keyring.ErrNotFound {
		return "", fmt.Errorf("secret not found for profile %q key %q", profile, key)
	}
	if err != nil {
		return "", fmt.Errorf("keyring get %q: %w", account, err)
	}
	return value, nil
}

// Delete removes a secret from the OS keyring for the given profile and key.
func (k *KeyringProvider) Delete(profile, key string) error {
	account := k.account(profile, key)
	if err := keyring.Delete(k.service, account); err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("keyring delete %q: %w", account, err)
	}
	return nil
}

// account returns the keyring account string combining profile and key.
func (k *KeyringProvider) account(profile, key string) string {
	return fmt.Sprintf("%s/%s", profile, key)
}
