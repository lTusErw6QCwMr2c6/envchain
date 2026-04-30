package secret

import (
	"fmt"
	"os"
)

// DefaultProvider returns a Provider based on the ENVCHAIN_PROVIDER environment
// variable. Supported values: env, keyring, vault, aws, doppler.
// Falls back to keyring if unset.
func DefaultProvider() (Provider, error) {
	providerType := ProviderType(os.Getenv("ENVCHAIN_PROVIDER"))
	if providerType == "" {
		providerType = ProviderKeyring
	}

	switch providerType {
	case ProviderEnv:
		return NewEnvProvider(), nil

	case ProviderKeyring:
		return NewKeyringProvider(), nil

	case ProviderVault:
		addr := os.Getenv("VAULT_ADDR")
		token := os.Getenv("VAULT_TOKEN")
		if addr == "" || token == "" {
			return nil, fmt.Errorf("vault provider requires VAULT_ADDR and VAULT_TOKEN")
		}
		return NewVaultProvider(addr, token), nil

	case ProviderAWS:
		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = "us-east-1"
		}
		return NewAWSProvider(region), nil

	case ProviderDoppler:
		token := os.Getenv("DOPPLER_TOKEN")
		project := os.Getenv("DOPPLER_PROJECT")
		config := os.Getenv("DOPPLER_CONFIG")
		if token == "" || project == "" || config == "" {
			return nil, fmt.Errorf("doppler provider requires DOPPLER_TOKEN, DOPPLER_PROJECT, and DOPPLER_CONFIG")
		}
		return NewDopplerProvider(token, project, config), nil

	default:
		return nil, fmt.Errorf("unknown provider type: %q (supported: env, keyring, vault, aws, doppler)", providerType)
	}
}

// AvailableProviders returns a list of all supported provider type strings.
func AvailableProviders() []ProviderType {
	return []ProviderType{
		ProviderEnv,
		ProviderKeyring,
		ProviderVault,
		ProviderAWS,
		ProviderDoppler,
	}
}
