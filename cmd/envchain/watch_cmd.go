package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	execpkg "github.com/envchain/envchain/internal/exec"
	"github.com/envchain/envchain/internal/secret"
)

func newWatchCmd() *cobra.Command {
	var interval time.Duration

	cmd := &cobra.Command{
		Use:   "watch <profile> -- <command> [args...]",
		Short: "Run a command and restart it when profile vars change",
		Long: `watch resolves the given profile (and its chain), injects the
environment variables into the command, and re-runs the command
whenever any variable in the chain changes.`,
		Args:               cobra.MinimumNArgs(2),
		DisableFlagParsing: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			cmdArgs := args[1:]

			store, err := defaultStore()
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			provider, err := secret.DefaultProvider()
			if err != nil {
				return fmt.Errorf("secret provider: %w", err)
			}

			wr := execpkg.NewWatchRunner(store, provider, interval)

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			fmt.Fprintf(os.Stderr, "[envchain] watching profile %q (poll every %s)\n", profileName, interval)
			return wr.Run(ctx, profileName, cmdArgs)
		},
	}

	cmd.Flags().DurationVarP(&interval, "interval", "i", 5*time.Second, "poll interval for profile changes")
	return cmd
}
