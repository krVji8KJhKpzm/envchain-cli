package env

import (
	"errors"
	"strings"
	"testing"
)

func suffixRotator(suffix string) RotateFunc {
	return func(project, name, oldValue string) (string, error) {
		return oldValue + suffix, nil
	}
}

func errorRotator(msg string) RotateFunc {
	return func(project, name, oldValue string) (string, error) {
		return "", errors.New(msg)
	}
}

func TestRotateVar(t *testing.T) {
	s := newMockStore()
	s.data["myapp"]["API_KEY"] = "old-secret"

	res, err := RotateVar(s, "myapp", "API_KEY", suffixRotator("-new"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.NewVal != "old-secret-new" {
		t.Errorf("expected 'old-secret-new', got %q", res.NewVal)
	}
	if !res.OldSet {
		t.Error("expected OldSet=true")
	}
	if res.Project != "myapp" || res.Var != "API_KEY" {
		t.Errorf("unexpected result metadata: %+v", res)
	}
	if s.data["myapp"]["API_KEY"] != "old-secret-new" {
		t.Error("store was not updated with new value")
	}
}

func TestRotateVarNotPreviouslySet(t *testing.T) {
	s := newMockStore()

	res, err := RotateVar(s, "myapp", "NEW_KEY", suffixRotator("generated"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OldSet {
		t.Error("expected OldSet=false for a new variable")
	}
	if res.NewVal != "generated" {
		t.Errorf("expected 'generated', got %q", res.NewVal)
	}
}

func TestRotateVarRotatorError(t *testing.T) {
	s := newMockStore()
	_, err := RotateVar(s, "myapp", "API_KEY", errorRotator("provider unavailable"))
	if err == nil {
		t.Fatal("expected error from rotator")
	}
	if !strings.Contains(err.Error(), "provider unavailable") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRotateVarInvalidProject(t *testing.T) {
	s := newMockStore()
	_, err := RotateVar(s, "", "API_KEY", suffixRotator("-new"))
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestRotateVarInvalidName(t *testing.T) {
	s := newMockStore()
	_, err := RotateVar(s, "myapp", "invalid-name!", suffixRotator("-new"))
	if err == nil {
		t.Fatal("expected error for invalid var name")
	}
}

func TestRotateAll(t *testing.T) {
	s := newMockStore()
	s.data["myapp"]["KEY_A"] = "aaa"
	s.data["myapp"]["KEY_B"] = "bbb"

	results, err := RotateAll(s, "myapp", []string{"KEY_A", "KEY_B"}, suffixRotator("_rotated"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].NewVal != "aaa_rotated" {
		t.Errorf("KEY_A: expected 'aaa_rotated', got %q", results[0].NewVal)
	}
	if results[1].NewVal != "bbb_rotated" {
		t.Errorf("KEY_B: expected 'bbb_rotated', got %q", results[1].NewVal)
	}
}
