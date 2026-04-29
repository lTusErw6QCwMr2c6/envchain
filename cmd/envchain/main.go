// Package main is the entry point for the envchain CLI tool.
// envchain manages and chains environment variable profiles across projects
// with optional secret store integration.
package main

import (
	"fmt"
	"os"

	"github.com/envchain/envchain/internal/exec"
	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"

	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "envchain",
		Short:   "Manage and chain environment variable profiles",
		Long:    `envchain lets you define, store, and inject environment variable profiles into processes, with optional integration into system secret stores.`,
		Version: version,
		SilenceUsage: true,
	}

	root.AddCommand(
		newRunCmd(),
		newSetCmd(),
		newListCmd(),
		newShowCmd(),
	)

	return root
}

// newRunCmd returns the "run" subcommand which executes a command with a
// profile's environment variables injected.
func newRunCmd() *cobra.Command {
	var profileNames []string

	cmd := &cobra.Command{
		Use:   "run -p <profile> [-- command args...]",
		Short: "Run a command with profile environment variables",
		Example: `  envchain run -p myapp -- ./server
  envchain run -p base -p override -- make deploy`,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := profile.DefaultStore()
			if err != nil {
				return fmt.Errorf("open profile store: %w", err)
			}
			provider := secret.DefaultProvider()
			runner := exec.NewRunner(store, provider)
			return runner.Run(cmd.Context(), profileNames, args)
		},
	}

	cmd.Flags().StringArrayVarP(&profileNames, "profile", "p", nil, "profile name(s) to load (can be repeated)")
	_ = cmd.MarkFlagRequired("profile")

	return cmd
}

// newSetCmd returns the "set" subcommand which creates or updates a profile.
func newSetCmd() *cobra.Command {
	var vars []string

	cmd := &cobra.Command{
		Use:   "set <profile>",
		Short: "Create or update a profile with environment variables",
		Args:  cobra.ExactArgs(1),
		Example: `  envchain set myapp -e PORT=8080 -e DEBUG=true`,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := profile.DefaultStore()
			if err != nil {
				return fmt.Errorf("open profile store: %w", err)
			}

			p := &profile.Profile{
				Name: args[0],
				Vars: vars,
			}
			if err := store.Save(p); err != nil {
				return fmt.Errorf("save profile: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Profile %q saved.\n", p.Name)
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&vars, "env", "e", nil, "environment variable in KEY=VALUE format (can be repeated)")
	return cmd
}

// newListCmd returns the "list" subcommand which prints all known profiles.
func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all stored profiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			store, err := profile.DefaultStore()
			if err != nil {
				return fmt.Errorf("open profile store: %w", err)
			}
			names, err := store.List()
			if err != nil {
				return fmt.Errorf("list profiles: %w", err)
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No profiles found.")
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}
}

// newShowCmd returns the "show" subcommand which prints the variables in a profile.
func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <profile>",
		Short: "Show environment variables defined in a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := profile.DefaultStore()
			if err != nil {
				return fmt.Errorf("open profile store: %w", err)
			}
			p, err := store.Load(args[0])
			if err != nil {
				return fmt.Errorf("load profile: %w", err)
			}
			for _, v := range p.Vars {
				fmt.Fprintln(cmd.OutOrStdout(), v)
			}
			return nil
		},
	}
}
