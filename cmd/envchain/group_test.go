package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findGroupSubCmd(name string) *cobra.Command {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "group" {
			for _, sub := range c.Commands() {
				if sub.Name() == name {
					return sub
				}
			}
		}
	}
	return nil
}

func TestGroupCmdRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "group" {
			return
		}
	}
	t.Fatal("group command not registered")
}

func TestGroupSubcommandsRegistered(t *testing.T) {
	for _, name := range []string{"add", "remove", "list"} {
		if findGroupSubCmd(name) == nil {
			t.Errorf("subcommand %q not registered under group", name)
		}
	}
}

func TestGroupAddRequiresArgs(t *testing.T) {
	cmd := findGroupSubCmd("add")
	if cmd == nil {
		t.Fatal("add subcommand not found")
	}
	if err := cmd.Args(cmd, []string{"proj", "grp"}); err == nil {
		t.Fatal("expected error with fewer than 3 args")
	}
}

func TestGroupRemoveRequiresArgs(t *testing.T) {
	cmd := findGroupSubCmd("remove")
	if cmd == nil {
		t.Fatal("remove subcommand not found")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Fatal("expected error with fewer than 3 args")
	}
}

func TestGroupListRequiresArgs(t *testing.T) {
	cmd := findGroupSubCmd("list")
	if cmd == nil {
		t.Fatal("list subcommand not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Fatal("expected error with zero args")
	}
	if err := cmd.Args(cmd, []string{"proj", "grp", "extra"}); err == nil {
		t.Fatal("expected error with more than 2 args")
	}
}
