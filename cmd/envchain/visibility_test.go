package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findVisibilitySubCmd(name string) *cobra.Command {
	for _, c := range rootCmd.Commands() {
		if c.Use == "visibility" {
			for _, sub := range c.Commands() {
				if sub.Use == name {
					return sub
				}
			}
		}
	}
	return nil
}

func TestVisibilityCmdRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "visibility" {
			return
		}
	}
	t.Error("visibility command not registered")
}

func TestVisibilitySubcommandsRegistered(t *testing.T) {
	for _, name := range []string{"set <project> <var> <level>", "get <project> <var>", "list <project> <level>", "remove <project> <var>"} {
		if findVisibilitySubCmd(name) == nil {
			t.Errorf("subcommand %q not registered", name)
		}
	}
}

func TestVisibilitySetRequiresArgs(t *testing.T) {
	cmd := findVisibilitySubCmd("set <project> <var> <level>")
	if cmd == nil {
		t.Fatal("set subcommand not found")
	}
	if err := cmd.Args(cmd, []string{"proj", "VAR"}); err == nil {
		t.Error("expected error when fewer than 3 args provided")
	}
}

func TestVisibilityGetRequiresArgs(t *testing.T) {
	cmd := findVisibilitySubCmd("get <project> <var>")
	if cmd == nil {
		t.Fatal("get subcommand not found")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error when fewer than 2 args provided")
	}
}

func TestVisibilityListRequiresArgs(t *testing.T) {
	cmd := findVisibilitySubCmd("list <project> <level>")
	if cmd == nil {
		t.Fatal("list subcommand not found")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error when fewer than 2 args provided")
	}
}

func TestVisibilityRemoveRequiresArgs(t *testing.T) {
	cmd := findVisibilitySubCmd("remove <project> <var>")
	if cmd == nil {
		t.Fatal("remove subcommand not found")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error when fewer than 2 args provided")
	}
}
