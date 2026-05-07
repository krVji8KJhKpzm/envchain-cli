package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func init() {
	lockCmd := &cobra.Command{
		Use:   "lock",
		Short: "Manage variable locks to prevent accidental modification",
	}

	lockSetCmd := &cobra.Command{
		Use:   "set <project> <var>",
		Short: "Lock a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runLockSet,
	}

	lockRemoveCmd := &cobra.Command{
		Use:   "remove <project> <var>",
		Short: "Unlock a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runLockRemove,
	}

	lockStatusCmd := &cobra.Command{
		Use:   "status <project> <var>",
		Short: "Show lock status of a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runLockStatus,
	}

	lockCmd.AddCommand(lockSetCmd, lockRemoveCmd, lockStatusCmd)
	rootCmd.AddCommand(lockCmd)
}

func newLockStore(project string) (*env.LockStore, error) {
	kc, err := keychain.New("envchain")
	if err != nil {
		return nil, fmt.Errorf("keychain init: %w", err)
	}
	store, err := env.New(kc)
	if err != nil {
		return nil, fmt.Errorf("store init: %w", err)
	}
	return env.NewLockStore(store, project)
}

func runLockSet(cmd *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	ls, err := newLockStore(project)
	if err != nil {
		return err
	}
	actor := currentActor()
	if err := ls.Lock(varName, actor); err != nil {
		return fmt.Errorf("lock %s: %w", varName, err)
	}
	fmt.Fprintf(os.Stdout, "Locked %s/%s (by %s)\n", project, varName, actor)
	return nil
}

func runLockRemove(cmd *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	ls, err := newLockStore(project)
	if err != nil {
		return err
	}
	if err := ls.Unlock(varName); err != nil {
		return fmt.Errorf("unlock %s: %w", varName, err)
	}
	fmt.Fprintf(os.Stdout, "Unlocked %s/%s\n", project, varName)
	return nil
}

func runLockStatus(cmd *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	ls, err := newLockStore(project)
	if err != nil {
		return err
	}
	locked, err := ls.IsLocked(varName)
	if err != nil {
		return fmt.Errorf("status %s: %w", varName, err)
	}
	if !locked {
		fmt.Fprintf(os.Stdout, "%s/%s is NOT locked\n", project, varName)
		return nil
	}
	entry, err := ls.GetLockEntry(varName)
	if err != nil {
		return fmt.Errorf("get lock entry: %w", err)
	}
	fmt.Fprintf(os.Stdout, "%s/%s is LOCKED by %s at %s\n",
		project, varName, entry.LockedBy, entry.LockedAt.Format("2006-01-02 15:04:05 UTC"))
	return nil
}
