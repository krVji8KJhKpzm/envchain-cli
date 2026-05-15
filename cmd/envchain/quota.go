package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func newQuotaStore(project string) (*env.QuotaStore, error) {
	kc, err := keychain.New()
	if err != nil {
		return nil, fmt.Errorf("open keychain: %w", err)
	}
	return env.NewQuotaStore(kc, project)
}

func init() {
	quotaCmd := &cobra.Command{
		Use:   "quota",
		Short: "Manage per-project variable count quotas",
	}

	setCmd := &cobra.Command{
		Use:   "set <project> <limit>",
		Short: "Set the maximum number of variables for a project",
		Args:  cobra.ExactArgs(2),
		RunE:  runQuotaSet,
	}

	getCmd := &cobra.Command{
		Use:   "get <project>",
		Short: "Show the quota limit for a project",
		Args:  cobra.ExactArgs(1),
		RunE:  runQuotaGet,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <project>",
		Short: "Remove the quota limit for a project (reverts to default)",
		Args:  cobra.ExactArgs(1),
		RunE:  runQuotaRemove,
	}

	quotaCmd.AddCommand(setCmd, getCmd, removeCmd)
	rootCmd.AddCommand(quotaCmd)
}

func runQuotaSet(cmd *cobra.Command, args []string) error {
	project := args[0]
	limit, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid limit %q: %w", args[1], err)
	}
	qs, err := newQuotaStore(project)
	if err != nil {
		return err
	}
	if err := qs.SetLimit(limit); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Quota for project %q set to %d variables.\n", project, limit)
	return nil
}

func runQuotaGet(cmd *cobra.Command, args []string) error {
	project := args[0]
	qs, err := newQuotaStore(project)
	if err != nil {
		return err
	}
	limit, err := qs.GetLimit()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Project %q quota limit: %d variables.\n", project, limit)
	return nil
}

func runQuotaRemove(cmd *cobra.Command, args []string) error {
	project := args[0]
	qs, err := newQuotaStore(project)
	if err != nil {
		return err
	}
	if err := qs.RemoveLimit(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Quota limit removed for project %q (default: %d).\n", project, 100)
	return nil
}
