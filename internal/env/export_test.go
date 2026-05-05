package env_test

import (
	"strings"
	"testing"

	"github.com/envchain-cli/internal/env"
)

func TestFormatShell(t *testing.T) {
	entries := []env.ExportEntry{
		{Name: "DB_URL", Value: "postgres://localhost/db"},
		{Name: "API_KEY", Value: "abc123"},
	}
	out := env.FormatShell(entries)
	if !strings.Contains(out, "export DB_URL=") {
		t.Errorf("expected DB_URL export, got:\n%s", out)
	}
	if !strings.Contains(out, "export API_KEY=") {
		t.Errorf("expected API_KEY export, got:\n%s", out)
	}
}

func TestFormatShellEmpty(t *testing.T) {
	out := env.FormatShell(nil)
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}

func TestFormatDotenv(t *testing.T) {
	entries := []env.ExportEntry{
		{Name: "PLAIN", Value: "simple"},
		{Name: "SPACED", Value: "hello world"},
		{Name: "QUOTED", Value: `say "hi"`},
	}
	out := env.FormatDotenv(entries)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "PLAIN=simple" {
		t.Errorf("unexpected line: %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], `SPACED="`) {
		t.Errorf("expected quoted value for SPACED, got %q", lines[1])
	}
	if !strings.Contains(lines[2], `\"`) {
		t.Errorf("expected escaped quote in QUOTED, got %q", lines[2])
	}
}
