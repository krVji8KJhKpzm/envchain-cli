package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findBookmarkSubCmd(name string) *cobra.Command {
	for _, c := range rootCmd.Commands() {
		if c.Use == "bookmark" {
			for _, sub := range c.Commands() {
				if sub.Use == name {
					return sub
				}
			}
		}
	}
	return nil
}

func TestBookmarkCmdRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "bookmark" {
			return
		}
	}
	t.Fatal("bookmark command not registered")
}

func TestBookmarkSubcommandsRegistered(t *testing.T) {
	for _, name := range []string{"set <name> <project> <variable>", "resolve <name>", "remove <name>", "list"} {
		if findBookmarkSubCmd(name) == nil {
			t.Errorf("subcommand %q not registered", name)
		}
	}
}

func TestBookmarkSetRequiresArgs(t *testing.T) {
	cmd := findBookmarkSubCmd("set <name> <project> <variable>")
	if cmd == nil {
		t.Fatal("set subcommand not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Fatal("expected error when no args provided")
	}
	if err := cmd.Args(cmd, []string{"name", "project"}); err == nil {
		t.Fatal("expected error when only two args provided")
	}
}

func TestBookmarkResolveRequiresArgs(t *testing.T) {
	cmd := findBookmarkSubCmd("resolve <name>")
	if cmd == nil {
		t.Fatal("resolve subcommand not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestBookmarkRemoveRequiresArgs(t *testing.T) {
	cmd := findBookmarkSubCmd("remove <name>")
	if cmd == nil {
		t.Fatal("remove subcommand not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Fatal("expected error when no args provided")
	}
}
