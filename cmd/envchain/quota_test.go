package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findQuotaSubCmd(name string) *cobra.Command {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "quota" {
			for _, child := range sub.Commands() {
				if child.Name() == name {
					return child
				}
			}
		}
	}
	return nil
}

func TestQuotaCmdRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "quota" {
			return
		}
	}
	t.Error("quota command not registered")
}

func TestQuotaSubcommandsRegistered(t *testing.T) {
	for _, name := range []string{"set", "get", "remove"} {
		if findQuotaSubCmd(name) == nil {
			t.Errorf("quota subcommand %q not registered", name)
		}
	}
}

func TestQuotaSetRequiresArgs(t *testing.T) {
	cmd := findQuotaSubCmd("set")
	if cmd == nil {
		t.Fatal("quota set command not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error with no args")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error with one arg")
	}
	if err := cmd.Args(cmd, []string{"proj", "10"}); err != nil {
		t.Errorf("unexpected error with two args: %v", err)
	}
}

func TestQuotaGetRequiresArgs(t *testing.T) {
	cmd := findQuotaSubCmd("get")
	if cmd == nil {
		t.Fatal("quota get command not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error with no args")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err != nil {
		t.Errorf("unexpected error with one arg: %v", err)
	}
}

func TestQuotaRemoveRequiresArgs(t *testing.T) {
	cmd := findQuotaSubCmd("remove")
	if cmd == nil {
		t.Fatal("quota remove command not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error with no args")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err != nil {
		t.Errorf("unexpected error with one arg: %v", err)
	}
}
