package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findPinSubCmd(name string) *cobra.Command {
	for _, c := range rootCmd.Commands() {
		if c.Use == "pin" {
			for _, sub := range c.Commands() {
				if sub.Use == name {
					return sub
				}
			}
		}
	}
	return nil
}

func TestPinCmdRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "pin" {
			return
		}
	}
	t.Error("expected 'pin' command to be registered")
}

func TestPinSubcommandsRegistered(t *testing.T) {
	for _, name := range []string{"set <project> <var>", "remove <project> <var>", "list <project>"} {
		if findPinSubCmd(name) == nil {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}

func TestPinSetRequiresArgs(t *testing.T) {
	cmd := findPinSubCmd("set <project> <var>")
	if cmd == nil {
		t.Fatal("pin set command not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error when no args provided to pin set")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error when only one arg provided to pin set")
	}
}

func TestPinRemoveRequiresArgs(t *testing.T) {
	cmd := findPinSubCmd("remove <project> <var>")
	if cmd == nil {
		t.Fatal("pin remove command not found")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error when only one arg provided to pin remove")
	}
}

func TestPinListRequiresArgs(t *testing.T) {
	cmd := findPinSubCmd("list <project>")
	if cmd == nil {
		t.Fatal("pin list command not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error when no args provided to pin list")
	}
}
