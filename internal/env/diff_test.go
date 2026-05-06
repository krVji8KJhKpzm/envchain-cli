package env_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain-cli/internal/env"
)

func TestDiffProjectsNoChanges(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("proj", "KEY", "val")
	_ = s.Set("proj2", "KEY", "val")

	diff, err := env.DiffProjects(s, "proj", "proj2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff.HasChanges() {
		t.Errorf("expected no changes, got: %s", diff)
	}
}

func TestDiffProjectsAdded(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("base", "A", "1")
	_ = s.Set("next", "A", "1")
	_ = s.Set("next", "B", "2")

	diff, err := env.DiffProjects(s, "base", "next")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.Added) != 1 || diff.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", diff.Added)
	}
}

func TestDiffProjectsRemoved(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("base", "A", "1")
	_ = s.Set("base", "B", "2")
	_ = s.Set("next", "A", "1")

	diff, err := env.DiffProjects(s, "base", "next")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.Removed) != 1 || diff.Removed[0] != "B" {
		t.Errorf("expected Removed=[B], got %v", diff.Removed)
	}
}

func TestDiffProjectsChanged(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("base", "KEY", "old")
	_ = s.Set("next", "KEY", "new")

	diff, err := env.DiffProjects(s, "base", "next")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.Changed) != 1 || diff.Changed[0] != "KEY" {
		t.Errorf("expected Changed=[KEY], got %v", diff.Changed)
	}
}

func TestDiffProjectsEmptyProject(t *testing.T) {
	s := newTestStore(t)
	_, err := env.DiffProjects(s, "", "next")
	if err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestDiffResultString(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("a", "ADDED", "1")
	_ = s.Set("b", "ADDED", "1")
	_ = s.Set("b", "EXTRA", "x")

	diff, _ := env.DiffProjects(s, "a", "b")
	out := diff.String()
	if !strings.Contains(out, "+ EXTRA") {
		t.Errorf("expected '+ EXTRA' in output, got: %s", out)
	}
}

func TestDiffResultStringNoChanges(t *testing.T) {
	d := env.DiffResult{}
	if d.String() != "no changes" {
		t.Errorf("expected 'no changes', got %q", d.String())
	}
}
