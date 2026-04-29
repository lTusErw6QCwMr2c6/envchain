// Package exec provides functionality for running commands with environment
// variables resolved from envchain profiles.
//
// The Runner type loads a named profile from the profile store, resolves any
// secret references through the configured secret.Provider, and executes the
// given command with the resulting environment variables merged into the
// current process environment.
//
// Example usage:
//
//	store, _ := profile.NewStore(profile.DefaultStore())
//	provider := secret.DefaultProvider()
//	runner := exec.NewRunner(store, provider)
//	runner.Run("myproject", []string{"go", "run", "."})
package exec
