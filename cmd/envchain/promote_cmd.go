package main

import (
	"fmt"

	"github.com/nicholasgasior/envchain/internal/profile"
	"github.com/spf13/cobra"
)

func newPromoteCmd() *cobra.Command {
	var overwrite bool
	var stripParents bool
	var tag string

	cmd := &cobra.Command{
		Use:   "promote <src> <dst>",
		Short: "Promote a profile from one environment to another",
		Long: `Copies a profile to a new name, optionally stripping parent references
and tagging the result. Useful for promoting staging configs to production.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := profile.DefaultStore()
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			src, dst := args[0], args[1]
			opts := profile.PromoteOptions{
				Overwrite:    overwrite,
				StripParents: stripParents,
				SuffixTag:    tag,
			}

			if err := profile.PromoteProfile(st, src, dst, opts); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "promoted %q → %q\n", src, dst)
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite destination profile if it exists")
	cmd.Flags().BoolVar(&stripParents, "strip-parents", false, "Remove parent references from the promoted profile")
	cmd.Flags().StringVar(&tag, "tag", "", "Tag to apply to the promoted profile")

	return cmd
}
