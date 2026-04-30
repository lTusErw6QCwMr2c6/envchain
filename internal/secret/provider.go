package secret

// Provider defines the interface for secret backends.
type Provider interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

// ProviderType identifies which backend to use.
type ProviderType string

const (
	ProviderEnv     ProviderType = "env"
	ProviderKeyring ProviderType = "keyring"
	ProviderVault   ProviderType = "vault"
	ProviderAWS     ProviderType = "aws"
	ProviderDoppler ProviderType = "doppler"
)

// ErrNotFound is returned when a secret key does not exist in the backend.
type ErrNotFound struct {
	Key string
}

func (e ErrNotFound) Error() string {
	return "secret not found: " + e.Key
}

func (e ErrNotFound) Is(target error) bool {
	_, ok := target.(ErrNotFound)
	return ok
}
