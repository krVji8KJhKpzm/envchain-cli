package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findWatchCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Name() == "watch" {
			return c
		}
	}
	return nil
}

func TestWatchCmdRegistered(t *testing.T) {
	if findWatchCmd(rootCmd) == nil {
		t.Fatal("watch command not registered")
	}
}

func TestWatchRequiresProjectAndVar(t *testing.T) {
	cmd := findWatchCmd(rootCmd)
	if cmd == nil {
		t.Fatal("watch command not found")
	}

	cmd.SetArgs([]string{})
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error with no args")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error with only project, no var")
	}
	if err := cmd.Args(cmd, []string{"proj", "VAR"}); err != nil {
		t.Errorf("unexpected error with valid args: %v", err)
	}
}

func TestWatchIntervalFlag(t *testing.T) {
	cmd := findWatchCmd(rootCmd)
	if cmd == nil {
		t.Fatal("watch command not found")
	}
	f := cmd.Flags().Lookup("interval")
	if f == nil {
		t.Fatal("--interval flag not registered")
	}
	if f.DefValue != "5s" {
		t.Errorf("expected default interval 5s, got %s", f.DefValue)
	}
}

func TestWatchShortIntervalFlag(t *testing.T) {
	cmd := findWatchCmd(rootCmd)
	if cmd == nil {
		t.Fatal("watch command not found")
	}
	f := cmd.Flags().ShorthandLookup("i")
	if f == nil {
		t.Fatal("-i shorthand not registered for --interval")
	}
}
