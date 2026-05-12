package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain/internal/profile"
)

func newAuditCmd() *cobra.Command {
	var filterProfile string
	var filterType string
	var last bool

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Show the audit log for profile changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := profile.DefaultStore()
			if err != nil {
				return err
			}
			s := profile.NewAuditStore(dir)
			log, err := s.Load()
			if err != nil {
				return fmt.Errorf("loading audit log: %w", err)
			}

			if last {
				e := log.Last()
				if e == nil {
					fmt.Fprintln(cmd.OutOrStdout(), "audit log is empty")
					return nil
				}
				printEvents(cmd, []profile.AuditEvent{*e})
				return nil
			}

			events := log.Events
			if filterProfile != "" {
				events = log.FilterByProfile(filterProfile)
			} else if filterType != "" {
				events = log.FilterByType(profile.AuditEventType(filterType))
			}

			if len(events) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no audit events found")
				return nil
			}
			printEvents(cmd, events)
			return nil
		},
	}

	cmd.Flags().StringVarP(&filterProfile, "profile", "p", "", "Filter events by profile name")
	cmd.Flags().StringVarP(&filterType, "type", "t", "", "Filter events by type (created|updated|deleted|copied|renamed)")
	cmd.Flags().BoolVar(&last, "last", false, "Show only the most recent event")
	return cmd
}

func printEvents(cmd *cobra.Command, events []profile.AuditEvent) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tTYPE\tPROFILE\tACTOR\tDETAIL")
	for _, e := range events {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.EventType,
			e.ProfileName,
			e.Actor,
			e.Detail,
		)
	}
	_ = w.Flush()
}
