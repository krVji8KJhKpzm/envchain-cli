package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

var execCmd = &cobra.Command{
	Use:   "exec <project> -- <command> [args...]",
	Short: "Run a command with project env vars injected",
	Long: `Execute a command with the environment variables for the given
project injected into the process environment.

Example:
  envchain exec myapp -- node server.js
  envchain exec myapp -- env | grep API`,
	Args:               cobra.MinimumNArgs(2),
	DisableFlagParsing: false,
	RunE:               runExec,
}

func init() {
	rootCmd.AddCommand(execCmd)
}

func runExec(cmd *cobra.Command, args []string) error {
	project := args[0]
	cmdArgs := args[1:]

	if len(cmdArgs) == 0 {
		return fmt.Errorf("no command specified after project name")
	}

	kc, err := keychain.New()
	if err != nil {
		return fmt.Errorf("failed to open keychain: %w", err)
	}

	store := env.New(kc)
	vars, err := store.List(project)
	if err != nil {
		return fmt.Errorf("failed to list vars for project %q: %w", project, err)
	}

	// Build environment: inherit current env, then overlay project vars
	environ := os.Environ()
	for _, v := range vars {
		val, err := store.Get(project, v.Name)
		if err != nil {
			return fmt.Errorf("failed to get var %q: %w", v.Name, err)
		}
		environ = append(environ, fmt.Sprintf("%s=%s", v.Name, val))
	}

	binary, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		return fmt.Errorf("command not found: %s", cmdArgs[0])
	}

	return syscall.Exec(binary, cmdArgs, environ)
}
