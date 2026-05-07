package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain/internal/env"
	"github.com/yourorg/envchain/internal/keychain"
)

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags on environment variables",
	}

	addCmd := &cobra.Command{
		Use:   "add <project> <var> <tag>",
		Short: "Add a tag to a variable",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagAdd(args[0], args[1], args[2])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <project> <var> <tag>",
		Short: "Remove a tag from a variable",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagRemove(args[0], args[1], args[2])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <project> <var>",
		Short: "List tags on a variable",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagList(args[0], args[1])
		},
	}

	findCmd := &cobra.Command{
		Use:   "find <project> <tag>",
		Short: "Find variables with a given tag",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagFind(args[0], args[1])
		},
	}

	tagCmd.AddCommand(addCmd, removeCmd, listCmd, findCmd)
	rootCmd.AddCommand(tagCmd)
}

func newTagStore() (*env.TagStore, error) {
	kc, err := keychain.New("envchain")
	if err != nil {
		return nil, err
	}
	s, err := env.New(kc)
	if err != nil {
		return nil, err
	}
	return env.NewTagStore(s)
}

func runTagAdd(project, varName, tag string) error {
	ts, err := newTagStore()
	if err != nil {
		return err
	}
	if err := ts.AddTag(project, varName, tag); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Tag %q added to %s/%s\n", tag, project, varName)
	return nil
}

func runTagRemove(project, varName, tag string) error {
	ts, err := newTagStore()
	if err != nil {
		return err
	}
	return ts.RemoveTag(project, varName, tag)
}

func runTagList(project, varName string) error {
	ts, err := newTagStore()
	if err != nil {
		return err
	}
	tags, err := ts.ListTags(project, varName)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		fmt.Println("(no tags)")
		return nil
	}
	fmt.Println(strings.Join(tags, "\n"))
	return nil
}

func runTagFind(project, tag string) error {
	ts, err := newTagStore()
	if err != nil {
		return err
	}
	names, err := ts.FindByTag(project, tag)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Println("(no variables found)")
		return nil
	}
	fmt.Println(strings.Join(names, "\n"))
	return nil
}
