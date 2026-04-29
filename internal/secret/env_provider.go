package secret

import (
	"fmt"
	"os"
	"strings"
)

const envPrefix = "ENVCHAIN_SECRET_"

// EnvProvider is a fallback secret provider that reads and writes secrets
// as environment variables prefixed with ENVCHAIN_SECRET_.
// Intended for testing and environments without a native keychain.
type EnvProvider struct{}

// NewEnvProvider creates a new EnvProvider.
func NewEnvProvider() *EnvProvider {
	return &EnvProvider{}
}

func envKey(service, key string) string {
	clean := func(s string) string {
		return strings.ToUpper(strings.NewReplacer("-", "_", ".", "_", " ", "_").Replace(s))
	}
	return fmt.Sprintf("%s%s__%s", envPrefix, clean(service), clean(key))
}

// Get retrieves a secret value from the environment.
func (p *EnvProvider) Get(service, key string) (string, error) {
	v, ok := os.LookupEnv(envKey(service, key))
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// Set stores a secret value in the current process environment.
func (p *EnvProvider) Set(service, key, value string) error {
	return os.Setenv(envKey(service, key), value)
}

// Delete removes a secret from the current process environment.
func (p *EnvProvider) Delete(service, key string) error {
	return os.Unsetenv(envKey(service, key))
}
