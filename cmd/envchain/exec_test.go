package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestExecInjectsEnvVars runs envchain exec via subprocess so we can verify
// that environment variables are actually injected into the child process.
func TestExecInjectsEnvVars(t *testing.T) {
	if os.Getenv("ENVCHAIN_INTEGRATION") == "" {
		t.Skip("skipping integration test; set ENVCHAIN_INTEGRATION=1 to run")
	}

	// Use the compiled binary; assumes `go build` has been run.
	binary := "./envchain"
	if _, err := os.Stat(binary); os.IsNotExist(err) {
		t.Skip("envchain binary not found; run `go build ./cmd/envchain` first")
	}

	// Set a known var in a test project.
	setCmd := exec.Command(binary, "set", "exectest", "EXEC_TEST_VAR=hello_world")
	if out, err := setCmd.CombinedOutput(); err != nil {
		t.Fatalf("set failed: %v\n%s", err, out)
	}
	t.Cleanup(func() {
		exec.Command(binary, "delete", "exectest", "EXEC_TEST_VAR").Run() //nolint:errcheck
	})

	// Execute `env` with the project and capture output.
	runCmd := exec.Command(binary, "exec", "exectest", "--", "env")
	out, err := runCmd.Output()
	if err != nil {
		t.Fatalf("exec failed: %v", err)
	}

	if !strings.Contains(string(out), "EXEC_TEST_VAR=hello_world") {
		t.Errorf("expected EXEC_TEST_VAR=hello_world in output, got:\n%s", out)
	}
}

func TestExecCmdRegistered(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "exec <project> -- <command> [args...]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("exec command not registered on rootCmd")
	}
}

func TestExecRequiresArgs(t *testing.T) {
	cmd := execCmd
	cmd.SetArgs([]string{})
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error when no args provided, got nil")
	}
}
