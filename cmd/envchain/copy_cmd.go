package main

import (
	"fmt"

	"github.com/nicholasgasior/envchain/internal/profile"
	"github.com/spf13/cobra"
)

func newCopyCmd(st *profile.Store) *cobra.Command {
	var overwrite bool
	var resetParents bool
	var rename bool

	cmd := &cobra.Command{
		Use:   "copy <src> <dst>",
		Short: "Copy or rename a profile",
		Long: `Copy duplicates an existing profile under a new name.
Use --rename to move (copy + delete) the source profile.
Use --reset-parents to clear the parent chain on the new profile.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]
			opts := profile.CopyOptions{
				Overwrite:    overwrite,
				ResetParents: resetParents,
			}

			var err error
			if rename {
				_, err = profile.RenameProfile(st, src, dst, opts)
				if err != nil {
					return fmt.Errorf("rename profile: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Profile %q renamed to %q\n", src, dst)
			} else {
				_, err = profile.CopyProfile(st, src, dst, opts)
				if err != nil {
					return fmt.Errorf("copy profile: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Profile %q copied to %q\n", src, dst)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite destination if it already exists")
	cmd.Flags().BoolVar(&resetParents, "reset-parents", false, "Clear parent chain on the copied profile")
	cmd.Flags().BoolVar(&rename, "rename", false, "Delete source after copying (rename semantics)")

	return cmd
}
