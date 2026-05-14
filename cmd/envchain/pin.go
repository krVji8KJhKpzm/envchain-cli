package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func newPinStore(project string) (*env.PinStore, error) {
	kc, err := keychain.New()
	if err != nil {
		return nil, err
	}
	return env.NewPinStore(kc, project)
}

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage pinned (write-protected) variables",
	}

	pinSetCmd := &cobra.Command{
		Use:   "set <project> <var>",
		Short: "Pin a variable to prevent modification",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPinSet(args[0], args[1])
		},
	}

	pinRemoveCmd := &cobra.Command{
		Use:   "remove <project> <var>",
		Short: "Unpin a variable",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPinRemove(args[0], args[1])
		},
	}

	pinListCmd := &cobra.Command{
		Use:   "list <project>",
		Short: "List all pinned variables in a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPinList(args[0])
		},
	}

	pinCmd.AddCommand(pinSetCmd, pinRemoveCmd, pinListCmd)
	rootCmd.AddCommand(pinCmd)
}

func runPinSet(project, varName string) error {
	ps, err := newPinStore(project)
	if err != nil {
		return err
	}
	if err := ps.Pin(varName); err != nil {
		return fmt.Errorf("pin: %w", err)
	}
	fmt.Fprintf(os.Stdout, "pinned %s in project %s\n", varName, project)
	return nil
}

func runPinRemove(project, varName string) error {
	ps, err := newPinStore(project)
	if err != nil {
		return err
	}
	if err := ps.Unpin(varName); err != nil {
		return fmt.Errorf("unpin: %w", err)
	}
	fmt.Fprintf(os.Stdout, "unpinned %s in project %s\n", varName, project)
	return nil
}

func runPinList(project string) error {
	ps, err := newPinStore(project)
	if err != nil {
		return err
	}
	pinned, err := ps.ListPinned()
	if err != nil {
		return fmt.Errorf("list pinned: %w", err)
	}
	if len(pinned) == 0 {
		fmt.Fprintln(os.Stdout, "no pinned variables")
		return nil
	}
	for _, v := range pinned {
		fmt.Fprintln(os.Stdout, v)
	}
	return nil
}
