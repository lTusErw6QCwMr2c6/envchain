package main

import (
	"fmt"
	"os"

	"github.com/envchain/envchain/internal/exec"
	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "export <profile>",
		Short: "Print environment variables for a profile to stdout",
		Long: `Resolve the full profile chain and print all environment variables
to stdout in the requested format (dotenv or export).

Example:
  envchain export myprofile
  envchain export --format=export myprofile | source /dev/stdin`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			store, err := profile.DefaultStore()
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			provider, err := secret.DefaultProvider()
			if err != nil {
				return fmt.Errorf("init secret provider: %w", err)
			}

			exporter := exec.NewExporter(store, provider)
			return exporter.Export(profileName, format, os.Stdout)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "dotenv",
		"Output format: dotenv or export")

	return cmd
}
