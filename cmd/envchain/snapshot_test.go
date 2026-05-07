package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findSnapshotSubCmd(name string) *cobra.Command {
	for _, c := range rootCmd.Commands() {
		if c.Use == "snapshot" {
			for _, sub := range c.Commands() {
				if strings.HasPrefix(sub.Use, name) {
					return sub
				}
			}
		}
	}
	return nil
}

func TestSnapshotCmdRegistered(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "snapshot" {
			found = true
			break
		}
	}
	if !found {
		t.Error("snapshot command not registered")
	}
}

func TestSnapshotSubcommandsRegistered(t *testing.T) {
	for _, name := range []string{"take", "list", "restore"} {
		if findSnapshotSubCmd(name) == nil {
			t.Errorf("snapshot sub-command %q not registered", name)
		}
	}
}

func TestSnapshotTakeRequiresProject(t *testing.T) {
	cmd := findSnapshotSubCmd("take")
	if cmd == nil {
		t.Fatal("take sub-command not found")
	}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error when project arg missing")
	}
}

func TestSnapshotRestoreRequiresTwoArgs(t *testing.T) {
	cmd := findSnapshotSubCmd("restore")
	if cmd == nil {
		t.Fatal("restore sub-command not found")
	}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	if err := cmd.RunE(cmd, []string{"myproject"}); err == nil {
		t.Error("expected error with only one arg")
	}
}

func TestSnapshotRestoreInvalidIndex(t *testing.T) {
	cmd := findSnapshotSubCmd("restore")
	if cmd == nil {
		t.Fatal("restore sub-command not found")
	}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.RunE(cmd, []string{"myproject", "notanumber"})
	if err == nil || !strings.Contains(err.Error(), "invalid index") {
		t.Errorf("expected invalid index error, got %v", err)
	}
}
