// Package exec provides utilities to execute commands with resolved
// environment profiles.
package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"
)

// WatchRunner re-runs a command whenever any of its source profiles change.
type WatchRunner struct {
	store    profile.Store
	provider secret.Provider
	interval time.Duration
	output   io.Writer
}

// NewWatchRunner creates a WatchRunner backed by the given store and provider.
func NewWatchRunner(store profile.Store, provider secret.Provider, interval time.Duration) *WatchRunner {
	return &WatchRunner{
		store:    store,
		provider: provider,
		interval: interval,
		output:   os.Stderr,
	}
}

// Run resolves profileName (and its chain), starts command args, and restarts
// it whenever a watched profile changes. It blocks until ctx is cancelled.
func (wr *WatchRunner) Run(ctx context.Context, profileName string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	names, err := profile.ChainNames(wr.store, profileName)
	if err != nil {
		return fmt.Errorf("resolve chain: %w", err)
	}

	watcher := profile.NewWatcher(wr.store, names, wr.interval)
	events := watcher.Watch(ctx)

	var cmd *exec.Cmd
	startCmd := func() error {
		env, err := wr.buildEnv(profileName)
		if err != nil {
			return err
		}
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Start()
	}

	if err := startCmd(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return nil
			}
			fmt.Fprintf(wr.output, "[envchain] profile %q %s — restarting\n", ev.ProfileName, ev.Kind)
			if cmd != nil && cmd.Process != nil {
				_ = cmd.Process.Kill()
				_ = cmd.Wait()
			}
			if err := startCmd(); err != nil {
				return fmt.Errorf("restart: %w", err)
			}
		case <-ctx.Done():
			if cmd != nil && cmd.Process != nil {
				_ = cmd.Process.Kill()
				_ = cmd.Wait()
			}
			return nil
		}
	}
}

func (wr *WatchRunner) buildEnv(profileName string) ([]string, error) {
	resolved, err := profile.ResolveChain(wr.store, profileName)
	if err != nil {
		return nil, err
	}
	return mergeEnv(os.Environ(), resolved.ToEnvMap()), nil
}
