package exec

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"
)

// Runner executes commands with environment variables resolved from profiles.
type Runner struct {
	store    *profile.Store
	provider secret.Provider
}

// NewRunner creates a new Runner with the given store and secret provider.
func NewRunner(store *profile.Store, provider secret.Provider) *Runner {
	return &Runner{
		store:    store,
		provider: provider,
	}
}

// Run executes the given command with environment variables from the named profile.
// The profile's variables are resolved via the secret provider and injected into
// the child process environment alongside the current process environment.
func (r *Runner) Run(profileName string, command []string) error {
	if len(command) == 0 {
		return fmt.Errorf("no command specified")
	}

	p, err := r.store.Load(profileName)
	if err != nil {
		return fmt.Errorf("load profile %q: %w", profileName, err)
	}

	envMap, err := p.ToEnvMap(r.provider)
	if err != nil {
		return fmt.Errorf("resolve secrets for profile %q: %w", profileName, err)
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = mergeEnv(os.Environ(), envMap)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}

// mergeEnv combines the base environment with overrides from envMap.
// Keys present in envMap take precedence over those in base.
func mergeEnv(base []string, overrides map[string]string) []string {
	result := make([]string, 0, len(base)+len(overrides))
	seen := make(map[string]bool, len(overrides))

	for k := range overrides {
		seen[k] = true
	}

	for _, entry := range base {
		if key := envKey(entry); !seen[key] {
			result = append(result, entry)
		}
	}

	for k, v := range overrides {
		result = append(result, k+"="+v)
	}
	return result
}

// envKey extracts the key portion from a KEY=VALUE string.
func envKey(entry string) string {
	for i := 0; i < len(entry); i++ {
		if entry[i] == '=' {
			return entry[:i]
		}
	}
	return entry
}
