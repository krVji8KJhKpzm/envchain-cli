package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func init() {
	expireCmd := &cobra.Command{
		Use:   "expire",
		Short: "Manage absolute expiration times for environment variables",
	}

	setCmd := &cobra.Command{
		Use:   "set <project> <var> <RFC3339-time>",
		Short: "Set an absolute expiry time for a variable",
		Args:  cobra.ExactArgs(3),
		RunE:  runExpireSet,
	}

	listCmd := &cobra.Command{
		Use:   "list <project> <var1> [var2...]",
		Short: "List expired variables from the given set",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runExpireList,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <project> <var>",
		Short: "Remove the expiry policy for a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runExpireRemove,
	}

	expireCmd.AddCommand(setCmd, listCmd, removeCmd)
	rootCmd.AddCommand(expireCmd)
}

func newExpiryStoreCmd(project string) (*env.ExpiryStore, error) {
	kc, err := keychain.New()
	if err != nil {
		return nil, fmt.Errorf("keychain init: %w", err)
	}
	return env.NewExpiryStore(kc, project)
}

func runExpireSet(_ *cobra.Command, args []string) error {
	project, varName, rawTime := args[0], args[1], args[2]
	at, err := time.Parse(time.RFC3339, rawTime)
	if err != nil {
		return fmt.Errorf("invalid time %q, expected RFC3339 (e.g. 2025-12-31T23:59:59Z): %w", rawTime, err)
	}
	s, err := newExpiryStoreCmd(project)
	if err != nil {
		return err
	}
	if err := s.SetExpiry(varName, at); err != nil {
		return fmt.Errorf("set expiry: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Expiry set for %s/%s: %s\n", project, varName, at.Format(time.RFC3339))
	return nil
}

func runExpireList(_ *cobra.Command, args []string) error {
	project := args[0]
	varNames := args[1:]
	s, err := newExpiryStoreCmd(project)
	if err != nil {
		return err
	}
	expired, err := s.ListExpired(varNames)
	if err != nil {
		return fmt.Errorf("list expired: %w", err)
	}
	if len(expired) == 0 {
		fmt.Fprintln(os.Stdout, "No expired variables found.")
		return nil
	}
	for _, e := range expired {
		fmt.Fprintln(os.Stdout, e.String())
	}
	return nil
}

func runExpireRemove(_ *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	s, err := newExpiryStoreCmd(project)
	if err != nil {
		return err
	}
	if err := s.RemoveExpiry(varName); err != nil {
		return fmt.Errorf("remove expiry: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Expiry removed for %s/%s\n", project, varName)
	return nil
}
