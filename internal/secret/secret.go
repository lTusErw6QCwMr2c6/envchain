package secret

import (
	"errors"
	"os"
	"runtime"
)

// Provider defines the interface for secret store backends.
type Provider interface {
	Get(service, key string) (string, error)
	Set(service, key, value string) error
	Delete(service, key string) error
}

// ErrNotFound is returned when a secret does not exist in the store.
var ErrNotFound = errors.New("secret not found")

// DefaultProvider returns the platform-appropriate secret provider.
// Falls back to an environment-variable-based provider if no native
// keychain is available.
func DefaultProvider() Provider {
	switch runtime.GOOS {
	case "darwin":
		return NewKeychainProvider()
	case "linux":
		if os.Getenv("ENVCHAIN_USE_KEYRING") == "1" {
			return NewKeychainProvider()
		}
		return NewEnvProvider()
	default:
		return NewEnvProvider()
	}
}
