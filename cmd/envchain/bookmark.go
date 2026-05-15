package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain-cli/internal/env"
	"github.com/envchain/envchain-cli/internal/keychain"
)

func newBookmarkStore() *env.BookmarkStore {
	kc := keychain.New()
	return env.NewBookmarkStore(kc)
}

func init() {
	bookmarkCmd := &cobra.Command{
		Use:   "bookmark",
		Short: "Manage named bookmarks to project variables",
	}

	setCmd := &cobra.Command{
		Use:   "set <name> <project> <variable>",
		Short: "Create or update a bookmark",
		Args:  cobra.ExactArgs(3),
		RunE:  runBookmarkSet,
	}

	resolveCmd := &cobra.Command{
		Use:   "resolve <name>",
		Short: "Print the project/variable a bookmark points to",
		Args:  cobra.ExactArgs(1),
		RunE:  runBookmarkResolve,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Delete a bookmark",
		Args:  cobra.ExactArgs(1),
		RunE:  runBookmarkRemove,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all bookmarks",
		Args:  cobra.NoArgs,
		RunE:  runBookmarkList,
	}

	bookmarkCmd.AddCommand(setCmd, resolveCmd, removeCmd, listCmd)
	rootCmd.AddCommand(bookmarkCmd)
}

func runBookmarkSet(cmd *cobra.Command, args []string) error {
	s := newBookmarkStore()
	if err := s.Set(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "bookmark %q set to %s/%s\n", args[0], args[1], args[2])
	return nil
}

func runBookmarkResolve(cmd *cobra.Command, args []string) error {
	s := newBookmarkStore()
	project, variable, err := s.Resolve(args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "%s/%s\n", project, variable)
	return nil
}

func runBookmarkRemove(cmd *cobra.Command, args []string) error {
	s := newBookmarkStore()
	if err := s.Remove(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "bookmark %q removed\n", args[0])
	return nil
}

func runBookmarkList(cmd *cobra.Command, args []string) error {
	s := newBookmarkStore()
	entries, err := s.List()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Fprintln(os.Stdout, "no bookmarks defined")
		return nil
	}
	for _, e := range entries {
		fmt.Fprintln(os.Stdout, e)
	}
	return nil
}
