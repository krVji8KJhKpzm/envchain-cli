package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func findSubCommand(root *cobra.Command, name string) *cobra.Command {
	for _, sub := range root.Commands() {
		if sub.Name() == name {
			return sub
		}
	}
	return nil
}

func TestTemplateCmdRegistered(t *testing.T) {
	cmd := findSubCommand(rootCmd, "template")
	if cmd == nil {
		t.Fatal("template command not registered")
	}
}

func TestTemplateRequiresProject(t *testing.T) {
	cmd := findSubCommand(rootCmd, "template")
	if cmd == nil {
		t.Skip("template command not found")
	}
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestTemplateListFlag(t *testing.T) {
	// Verify --list flag exists on the template command
	cmd := findSubCommand(rootCmd, "template")
	if cmd == nil {
		t.Skip("template command not found")
	}
	f := cmd.Flags().Lookup("list")
	if f == nil {
		t.Fatal("expected --list flag")
	}
}

func TestTemplateFileFlag(t *testing.T) {
	cmd := findSubCommand(rootCmd, "template")
	if cmd == nil {
		t.Skip("template command not found")
	}
	f := cmd.Flags().Lookup("file")
	if f == nil {
		t.Fatal("expected --file flag")
	}
}

func TestTemplateListOnlyOutput(t *testing.T) {
	cmd := findSubCommand(rootCmd, "template")
	if cmd == nil {
		t.Skip("template command not found")
	}
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.Flags().Set("list", "true")
	// runTemplate with list flag and no real keychain needed
	err := runTemplate(cmd, []string{"proj", "host={{DB_HOST}} port={{DB_PORT}}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Error("expected placeholder names in output")
	}
	_ = cmd.Flags().Set("list", "false")
}
