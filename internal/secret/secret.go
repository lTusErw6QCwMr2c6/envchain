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

// IsNotFound reports whether err is an ErrNotFound error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// DefaultProvider returns the platform-appropriate secret provider.
// Falls back to an environment-variable-based provider if no native
// keychain is available.
//
// On Linux, set ENVCHAIN_USE_KEYRING=1 to use the system keyring
// instead of the default environment-variable-based provider.
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
