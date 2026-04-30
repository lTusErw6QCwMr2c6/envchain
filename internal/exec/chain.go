// Package exec provides functionality for running commands with environment
// variables loaded from envchain profiles.
package exec

import (
	"fmt"

	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"
)

// ChainRunner resolves multiple profiles in order and merges their environment
// variables, with later profiles taking precedence over earlier ones.
type ChainRunner struct {
	store    *profile.Store
	provider secret.Provider
}

// NewChainRunner creates a new ChainRunner with the given store and provider.
func NewChainRunner(store *profile.Store, provider secret.Provider) *ChainRunner {
	return &ChainRunner{
		store:    store,
		provider: provider,
	}
}

// ResolveChain loads and merges environment variables from a list of profile
// names in order. Later profiles override earlier ones on key conflicts.
func (c *ChainRunner) ResolveChain(profileNames []string) (map[string]string, error) {
	if len(profileNames) == 0 {
		return nil, fmt.Errorf("no profiles specified")
	}

	merged := make(map[string]string)

	for _, name := range profileNames {
		p, err := c.store.Load(name)
		if err != nil {
			return nil, fmt.Errorf("loading profile %q: %w", name, err)
		}

		for _, v := range p.Vars {
			val, err := c.provider.Get(name, v.Name)
			if err != nil {
				return nil, fmt.Errorf("resolving secret %q in profile %q: %w", v.Name, name, err)
			}
			merged[v.Name] = val
		}
	}

	return merged, nil
}

// Run executes a command with the merged environment from the given profiles.
func (c *ChainRunner) Run(profileNames []string, command string, args []string) error {
	if command == "" {
		return fmt.Errorf("no command specified")
	}

	env, err := c.ResolveChain(profileNames)
	if err != nil {
		return err
	}

	r := NewRunner(c.store, c.provider)
	return r.RunWithEnv(env, command, args)
}
