package main

import (
	"fmt"
	"os/user"

	"github.com/nicholasgasior/envchain/internal/profile"
	"github.com/spf13/cobra"
)

func newPinCmd(store *profile.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage profile pins (snapshots at a point in time)",
	}
	cmd.AddCommand(newPinCreateCmd(store))
	cmd.AddCommand(newPinDiffCmd(store))
	return cmd
}

func newPinCreateCmd(store *profile.Store) *cobra.Command {
	var label string
	cmd := &cobra.Command{
		Use:   "create <profile>",
		Short: "Pin the current state of a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			who := currentUser()
			pin, err := profile.PinProfile(store, name, who)
			if err != nil {
				return err
			}
			dir := store.Dir()
			ps, err := profile.NewPinStore(dir)
			if err != nil {
				return err
			}
			if err := ps.Save(label, pin); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pinned %q as %q (by %s)\n", name, label, who)
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "default", "Label for this pin")
	return cmd
}

func newPinDiffCmd(store *profile.Store) *cobra.Command {
	var label string
	cmd := &cobra.Command{
		Use:   "diff <profile>",
		Short: "Show changes since a pin was taken",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			dir := store.Dir()
			ps, err := profile.NewPinStore(dir)
			if err != nil {
				return err
			}
			pin, err := ps.Load(name, label)
			if err != nil {
				return err
			}
			diff, err := profile.DiffPin(store, pin)
			if err != nil {
				return err
			}
			printDiff(cmd, diff)
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "default", "Label of pin to diff against")
	return cmd
}

func currentUser() string {
	u, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return u.Username
}

func printDiff(cmd *cobra.Command, d profile.DiffResult) {
	w := cmd.OutOrStdout()
	for _, k := range d.Added {
		fmt.Fprintf(w, "+ %s\n", k)
	}
	for _, k := range d.Removed {
		fmt.Fprintf(w, "- %s\n", k)
	}
	for _, k := range d.Changed {
		fmt.Fprintf(w, "~ %s\n", k)
	}
	if len(d.Added)+len(d.Removed)+len(d.Changed) == 0 {
		fmt.Fprintln(w, "No changes since pin.")
	}
}
