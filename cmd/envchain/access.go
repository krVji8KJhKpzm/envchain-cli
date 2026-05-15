package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func init() {
	accessCmd := &cobra.Command{
		Use:   "access",
		Short: "Manage variable access logs",
	}

	recordCmd := &cobra.Command{
		Use:   "record <project> <var>",
		Short: "Record a read access for a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runAccessRecord,
	}

	listCmd := &cobra.Command{
		Use:   "list <project> <var>",
		Short: "List access history for a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runAccessList,
	}

	clearCmd := &cobra.Command{
		Use:   "clear <project> <var>",
		Short: "Clear access history for a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runAccessClear,
	}

	accessCmd.AddCommand(recordCmd, listCmd, clearCmd)
	rootCmd.AddCommand(accessCmd)
}

func newAccessStore() *env.AccessStore {
	kc := keychain.New("envchain")
	return env.NewAccessStore(kc)
}

func runAccessRecord(cmd *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	actor := currentActor()
	s := newAccessStore()
	if err := s.Record(project, varName, actor); err != nil {
		return fmt.Errorf("record access: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Access recorded for %s/%s by %s\n", project, varName, actor)
	return nil
}

func runAccessList(cmd *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	s := newAccessStore()
	entries, err := s.List(project, varName)
	if err != nil {
		return fmt.Errorf("list access: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No access records found.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ACTOR\tVAR\tACCESSED AT")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.Actor, e.VarName, e.AccessedAt.Format("2006-01-02 15:04:05 UTC"))
	}
	return w.Flush()
}

func runAccessClear(cmd *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	s := newAccessStore()
	if err := s.Clear(project, varName); err != nil {
		return fmt.Errorf("clear access: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Access log cleared for %s/%s\n", project, varName)
	return nil
}
