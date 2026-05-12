package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nicholasgasior/envchain/internal/profile"
	"github.com/spf13/cobra"
)

func newLockCmd(store profile.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Manage profile locks",
	}
	cmd.AddCommand(newLockAcquireCmd(store))
	cmd.AddCommand(newLockReleaseCmd(store))
	cmd.AddCommand(newLockStatusCmd(store))
	return cmd
}

func newLockAcquireCmd(store profile.Store) *cobra.Command {
	var owner string
	var ttl time.Duration

	cmd := &cobra.Command{
		Use:   "acquire <profile>",
		Short: "Acquire a lock on a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ls := profile.NewLockStore(store)
			if err := ls.Acquire(args[0], owner, ttl); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Lock acquired on profile %q by %q (TTL: %s)\n", args[0], owner, ttl)
			return nil
		},
	}
	cmd.Flags().StringVar(&owner, "owner", currentLockUser(), "lock owner identifier")
	cmd.Flags().DurationVar(&ttl, "ttl", 30*time.Minute, "lock time-to-live duration")
	return cmd
}

func newLockReleaseCmd(store profile.Store) *cobra.Command {
	var owner string

	cmd := &cobra.Command{
		Use:   "release <profile>",
		Short: "Release a lock on a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ls := profile.NewLockStore(store)
			if err := ls.Release(args[0], owner); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Lock released on profile %q\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&owner, "owner", currentLockUser(), "lock owner identifier")
	return cmd
}

func newLockStatusCmd(store profile.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "status <profile>",
		Short: "Show lock status for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ls := profile.NewLockStore(store)
			entry, err := ls.Get(args[0])
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Profile %q is not locked\n", args[0])
				return nil
			}
			if entry.IsExpired() {
				fmt.Fprintf(cmd.OutOrStdout(), "Profile %q has an expired lock (owner: %q, expired: %s)\n",
					args[0], entry.Owner, entry.ExpiresAt.Format(time.RFC3339))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Profile %q is locked by %q (expires: %s)\n",
					args[0], entry.Owner, entry.ExpiresAt.Format(time.RFC3339))
			}
			return nil
		},
	}
}

func currentLockUser() string {
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	return "unknown"
}
