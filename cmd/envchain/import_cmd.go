package main

import (
	"fmt"
	"os"

	"github.com/envchain/envchain/internal/profile"
	"github.com/spf13/cobra"
)

func newImportCmd(st profile.Store) *cobra.Command {
	var (
		format    string
		overwrite bool
		file      string
	)

	cmd := &cobra.Command{
		Use:   "import <profile>",
		Short: "Import environment variables from a file into a profile",
		Long: `Read key=value pairs from a dotenv or export-style file and merge
them into the named profile. Existing variables are preserved unless
--overwrite is specified.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			var fmt_ profile.ImportFormat
			switch format {
			case "dotenv":
				fmt_ = profile.ImportFormatDotenv
			case "export":
				fmt_ = profile.ImportFormatExport
			default:
				return fmt.Errorf("unknown format %q: use 'dotenv' or 'export'", format)
			}

			r := cmd.InOrStdin()
			if file != "" && file != "-" {
				f, err := os.Open(file)
				if err != nil {
					return fmt.Errorf("open %q: %w", file, err)
				}
				defer f.Close()
				r = f
			}

			opts := profile.ImportOptions{Overwrite: overwrite}
			if err := profile.ImportProfile(st, name, r, fmt_, opts); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "imported variables into profile %q\n", name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "input format: dotenv or export")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing variables")
	cmd.Flags().StringVar(&file, "file", "-", "input file (default: stdin)")
	return cmd
}
