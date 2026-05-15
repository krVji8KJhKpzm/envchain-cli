package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findNoteSubCmd(name string) *cobra.Command {
	for _, c := range rootCmd.Commands() {
		if c.Use == "note" {
			for _, sub := range c.Commands() {
				if sub.Use == name {
					return sub
				}
			}
		}
	}
	return nil
}

func TestNoteCmdRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "note" {
			return
		}
	}
	t.Error("expected 'note' command to be registered")
}

func TestNoteSubcommandsRegistered(t *testing.T) {
	for _, sub := range []string{"set <project> <var> <note>", "get <project> <var>", "remove <project> <var>"} {
		if findNoteSubCmd(sub) == nil {
			t.Errorf("expected subcommand %q to be registered", sub)
		}
	}
}

func TestNoteSetRequiresArgs(t *testing.T) {
	cmd := findNoteSubCmd("set <project> <var> <note>")
	if cmd == nil {
		t.Fatal("note set command not found")
	}
	if err := cmd.Args(cmd, []string{"proj", "VAR"}); err == nil {
		t.Error("expected error when fewer than 3 args provided")
	}
}

func TestNoteGetRequiresArgs(t *testing.T) {
	cmd := findNoteSubCmd("get <project> <var>")
	if cmd == nil {
		t.Fatal("note get command not found")
	}
	if err := cmd.Args(cmd, []string{"proj"}); err == nil {
		t.Error("expected error when fewer than 2 args provided")
	}
}

func TestNoteRemoveRequiresArgs(t *testing.T) {
	cmd := findNoteSubCmd("remove <project> <var>")
	if cmd == nil {
		t.Fatal("note remove command not found")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error when no args provided")
	}
}
