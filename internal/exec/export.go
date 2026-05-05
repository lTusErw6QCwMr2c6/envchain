// Package exec provides utilities for running commands with injected
// environment variables sourced from envchain profiles.
package exec

import (
	"fmt"
	"io"

	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"
)

// Exporter resolves a profile chain and writes environment variable
// declarations to a writer in the requested format.
type Exporter struct {
	store    *profile.Store
	provider secret.Provider
}

// NewExporter creates an Exporter backed by the given store and secret provider.
func NewExporter(store *profile.Store, provider secret.Provider) *Exporter {
	return &Exporter{store: store, provider: provider}
}

// Export writes the resolved environment variables for the named profile
// (and its chain) to w using the specified format (dotenv or export).
func (e *Exporter) Export(profileName, format string, w io.Writer) error {
	names, err := profile.ChainNames(e.store, profileName)
	if err != nil {
		return fmt.Errorf("resolve chain: %w", err)
	}

	combined := map[string]string{}
	for _, name := range names {
		p, err := e.store.Load(name)
		if err != nil {
			return fmt.Errorf("load profile %q: %w", name, err)
		}
		expanded := profile.ExpandVars(p, combined)
		for _, v := range expanded.Vars {
			val, err := resolveValue(e.provider, name, v)
			if err != nil {
				return err
			}
			combined[v.Key] = val
		}
	}

	synthetic := &profile.Profile{Name: profileName}
	for k, v := range combined {
		synthetic.Vars = append(synthetic.Vars, profile.Var{Key: k, Value: v})
	}

	out, err := profile.RenderTemplate(synthetic, format)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(w, out)
	return err
}

func resolveValue(provider secret.Provider, profileName string, v profile.Var) (string, error) {
	if v.Secret {
		val, err := provider.Get(profileName, v.Key)
		if err != nil {
			return "", fmt.Errorf("get secret %q/%q: %w", profileName, v.Key, err)
		}
		return val, nil
	}
	return v.Value, nil
}
