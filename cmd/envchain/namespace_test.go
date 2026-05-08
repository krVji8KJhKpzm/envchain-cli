package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findNamespaceSubCmd(name string) *cobra.Command {
	for _, sub := range namespaceCmd.Commands() {
		if sub.Name() == name {
			return sub
		}
	}
	return nil
}

func TestNamespaceCmdRegistered(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Name() == "namespace" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("namespace command not registered on rootCmd")
	}
}

func TestNamespaceSubcommandsRegistered(t *testing.T) {
	for _, name := range []string{"list", "parse"} {
		if findNamespaceSubCmd(name) == nil {
			t.Errorf("subcommand %q not registered under namespace", name)
		}
	}
}

func TestNamespaceListRequiresArgs(t *testing.T) {
	cmd := findNamespaceSubCmd("list")
	if cmd == nil {
		t.Fatal("list subcommand not found")
	}
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Fatal("expected error when no args provided to list")
	}
}

func TestNamespaceParseRequiresArgs(t *testing.T) {
	cmd := findNamespaceSubCmd("parse")
	if cmd == nil {
		t.Fatal("parse subcommand not found")
	}
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Fatal("expected error when no args provided to parse")
	}
}

func TestNamespaceParseExactlyOneArg(t *testing.T) {
	cmd := findNamespaceSubCmd("parse")
	if cmd == nil {
		t.Fatal("parse subcommand not found")
	}
	err := cmd.Args(cmd, []string{"acme/payments", "extra"})
	if err == nil {
		t.Fatal("expected error when more than one arg provided")
	}
}
