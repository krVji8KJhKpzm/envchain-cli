package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func newDepStore() *env.DependencyStore {
	kc := keychain.New(serviceName)
	s := env.New(kc)
	return env.NewDependencyStore(s)
}

func init() {
	depCmd := &cobra.Command{
		Use:   "dep",
		Short: "Manage variable dependencies",
	}

	setCmd := &cobra.Command{
		Use:   "set <project> <var> <dep1,dep2,...>",
		Short: "Declare dependencies for a variable",
		Args:  cobra.ExactArgs(3),
		RunE:  runDepSet,
	}

	getCmd := &cobra.Command{
		Use:   "get <project> <var>",
		Short: "List declared dependencies for a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runDepGet,
	}

	checkCmd := &cobra.Command{
		Use:   "check <project> <var>",
		Short: "Check whether all dependencies are satisfied",
		Args:  cobra.ExactArgs(2),
		RunE:  runDepCheck,
	}

	depCmd.AddCommand(setCmd, getCmd, checkCmd)
	rootCmd.AddCommand(depCmd)
}

func runDepSet(cmd *cobra.Command, args []string) error {
	project, varName, rawDeps := args[0], args[1], args[2]
	deps := strings.Split(rawDeps, ",")
	ds := newDepStore()
	if err := ds.SetDeps(project, varName, deps); err != nil {
		return fmt.Errorf("dep set: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Dependencies set for %s/%s\n", project, varName)
	return nil
}

func runDepGet(cmd *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	ds := newDepStore()
	deps, err := ds.GetDeps(project, varName)
	if err != nil {
		return fmt.Errorf("dep get: %w", err)
	}
	if len(deps) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "(no dependencies declared)")
		return nil
	}
	for _, d := range deps {
		fmt.Fprintln(cmd.OutOrStdout(), d)
	}
	return nil
}

func runDepCheck(cmd *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	ds := newDepStore()
	ok, missing, err := ds.Satisfied(project, varName)
	if err != nil {
		return fmt.Errorf("dep check: %w", err)
	}
	if ok {
		fmt.Fprintf(cmd.OutOrStdout(), "All dependencies satisfied for %s/%s\n", project, varName)
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Missing dependencies for %s/%s:\n", project, varName)
	for _, m := range missing {
		fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", m)
	}
	return fmt.Errorf("unsatisfied dependencies: %s", strings.Join(missing, ", "))
}
