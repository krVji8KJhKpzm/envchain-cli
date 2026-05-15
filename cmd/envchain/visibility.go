package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain-cli/internal/env"
	"github.com/envchain/envchain-cli/internal/keychain"
)

func newVisibilityStore() *env.VisibilityStore {
	kc := keychain.New()
	s := env.New(kc)
	return env.NewVisibilityStore(s)
}

func init() {
	visibilityCmd := &cobra.Command{
		Use:   "visibility",
		Short: "Manage variable visibility levels (public, private, secret)",
	}

	setCmd := &cobra.Command{
		Use:   "set <project> <var> <level>",
		Short: "Set visibility level for a variable",
		Args:  cobra.ExactArgs(3),
		RunE:  runVisibilitySet,
	}

	getCmd := &cobra.Command{
		Use:   "get <project> <var>",
		Short: "Get visibility level of a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runVisibilityGet,
	}

	listCmd := &cobra.Command{
		Use:   "list <project> <level>",
		Short: "List variables with a given visibility level",
		Args:  cobra.ExactArgs(2),
		RunE:  runVisibilityList,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <project> <var>",
		Short: "Remove visibility setting for a variable (resets to default)",
		Args:  cobra.ExactArgs(2),
		RunE:  runVisibilityRemove,
	}

	visibilityCmd.AddCommand(setCmd, getCmd, listCmd, removeCmd)
	rootCmd.AddCommand(visibilityCmd)
}

func runVisibilitySet(cmd *cobra.Command, args []string) error {
	vs := newVisibilityStore()
	if err := vs.Set(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "visibility of %s/%s set to %s\n", args[0], args[1], strings.ToLower(args[2]))
	return nil
}

func runVisibilityGet(cmd *cobra.Command, args []string) error {
	vs := newVisibilityStore()
	level, err := vs.Get(args[0], args[1])
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, level)
	return nil
}

func runVisibilityList(cmd *cobra.Command, args []string) error {
	vs := newVisibilityStore()
	names, err := vs.ListByLevel(args[0], args[1])
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintf(os.Stdout, "no variables with visibility %q in project %q\n", args[1], args[0])
		return nil
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintln(os.Stdout, n)
	}
	return nil
}

func runVisibilityRemove(cmd *cobra.Command, args []string) error {
	vs := newVisibilityStore()
	if err := vs.Remove(args[0], args[1]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "visibility setting removed for %s/%s\n", args[0], args[1])
	return nil
}
