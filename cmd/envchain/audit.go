package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain-cli/internal/env"
)

// globalAuditLog is a package-level audit log shared across commands.
var globalAuditLog = env.NewAuditLog()

// currentActor returns the current OS username for audit purposes.
func currentActor() string {
	u, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return u.Username
}

func init() {
	var project string
	var action string

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Show audit log of variable operations",
		Long:  "Display a history of environment variable operations recorded during this session.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudit(project, action)
		},
	}

	auditCmd.Flags().StringVarP(&project, "project", "p", "", "Filter by project name")
	auditCmd.Flags().StringVarP(&action, "action", "a", "", "Filter by action (set, delete, rotate, copy, move)")

	rootCmd.AddCommand(auditCmd)
}

func runAudit(project, action string) error {
	events := globalAuditLog.Events()

	if project != "" {
		events = globalAuditLog.FilterByProject(project)
	}

	if action != "" {
		filtered := make([]env.AuditEvent, 0, len(events))
		for _, e := range events {
			if e.Action == action {
				filtered = append(filtered, e)
			}
		}
		events = filtered
	}

	if len(events) == 0 {
		fmt.Fprintln(os.Stderr, "No audit events found.")
		return nil
	}

	for _, e := range events {
		fmt.Println(e.String())
	}
	return nil
}
