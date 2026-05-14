package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain/internal/env"
	"github.com/yourorg/envchain/internal/keychain"
)

func newGroupStore() (*env.GroupStore, error) {
	kc, err := keychain.New()
	if err != nil {
		return nil, err
	}
	store, err := env.New(kc)
	if err != nil {
		return nil, err
	}
	return env.NewGroupStore(store)
}

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage variable groups within a project",
	}

	addCmd := &cobra.Command{
		Use:   "add <project> <group> <var>",
		Short: "Add a variable to a group",
		Args:  cobra.ExactArgs(3),
		RunE:  runGroupAdd,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <project> <group> <var>",
		Short: "Remove a variable from a group",
		Args:  cobra.ExactArgs(3),
		RunE:  runGroupRemove,
	}

	listCmd := &cobra.Command{
		Use:   "list <project> [group]",
		Short: "List groups or members of a group",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runGroupList,
	}

	groupCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(groupCmd)
}

func runGroupAdd(cmd *cobra.Command, args []string) error {
	gs, err := newGroupStore()
	if err != nil {
		return err
	}
	if err := gs.AddToGroup(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Added %q to group %q in project %q\n", args[2], args[1], args[0])
	return nil
}

func runGroupRemove(cmd *cobra.Command, args []string) error {
	gs, err := newGroupStore()
	if err != nil {
		return err
	}
	if err := gs.RemoveFromGroup(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Removed %q from group %q in project %q\n", args[2], args[1], args[0])
	return nil
}

func runGroupList(cmd *cobra.Command, args []string) error {
	gs, err := newGroupStore()
	if err != nil {
		return err
	}
	project := args[0]
	if len(args) == 1 {
		groups, err := gs.ListGroups(project)
		if err != nil {
			return err
		}
		if len(groups) == 0 {
			fmt.Fprintln(os.Stdout, "No groups defined.")
			return nil
		}
		fmt.Fprintln(os.Stdout, strings.Join(groups, "\n"))
		return nil
	}
	members, err := gs.ListGroup(project, args[1])
	if err != nil {
		return err
	}
	if len(members) == 0 {
		fmt.Fprintln(os.Stdout, "Group is empty.")
		return nil
	}
	fmt.Fprintln(os.Stdout, strings.Join(members, "\n"))
	return nil
}
