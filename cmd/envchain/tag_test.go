package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findTagSubCmd(t *testing.T, name string) *cobra.Command {
	t.Helper()
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "tag" {
			for _, s := range sub.Commands() {
				if s.Name() == name {
					return s
				}
			}
		}
	}
	t.Fatalf("tag subcommand %q not found", name)
	return nil
}

func TestTagCmdRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "tag" {
			return
		}
	}
	t.Error("expected 'tag' command to be registered")
}

func TestTagSubcommandsRegistered(t *testing.T) {
	for _, name := range []string{"add", "remove", "list", "find"} {
		findTagSubCmd(t, name)
	}
}

func TestTagAddRequiresArgs(t *testing.T) {
	cmd := findTagSubCmd(t, "add")
	if err := cmd.Args(cmd, []string{"proj", "VAR"}); err == nil {
		t.Error("expected error with fewer than 3 args")
	}
}

func TestTagFindRequiresArgs(t *testing.T) {
	cmd := findTagSubCmd(t, "find")
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error with fewer than 2 args")
	}
}

func TestTagListRequiresArgs(t *testing.T) {
	cmd := findTagSubCmd(t, "list")
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error with fewer than 2 args")
	}
}
