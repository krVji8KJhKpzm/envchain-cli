package env

import (
	"testing"
)

func newVisibilityStore(t *testing.T) *VisibilityStore {
	t.Helper()
	ms := newMockStore()
	return NewVisibilityStore(ms)
}

func TestVisibilitySetAndGet(t *testing.T) {
	vs := newVisibilityStore(t)
	if err := vs.Set("proj", "API_KEY", "secret"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	level, err := vs.Get("proj", "API_KEY")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if level != "secret" {
		t.Errorf("expected secret, got %s", level)
	}
}

func TestVisibilityDefaultPrivate(t *testing.T) {
	vs := newVisibilityStore(t)
	level, err := vs.Get("proj", "UNSET_VAR")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if level != VisibilityPrivate {
		t.Errorf("expected default private, got %s", level)
	}
}

func TestVisibilityInvalidLevel(t *testing.T) {
	vs := newVisibilityStore(t)
	if err := vs.Set("proj", "VAR", "classified"); err == nil {
		t.Error("expected error for invalid visibility level")
	}
}

func TestVisibilityRemove(t *testing.T) {
	vs := newVisibilityStore(t)
	_ = vs.Set("proj", "DB_PASS", "secret")
	if err := vs.Remove("proj", "DB_PASS"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	level, _ := vs.Get("proj", "DB_PASS")
	if level != VisibilityPrivate {
		t.Errorf("expected default after removal, got %s", level)
	}
}

func TestVisibilityListByLevel(t *testing.T) {
	vs := newVisibilityStore(t)
	_ = vs.Set("proj", "API_KEY", "secret")
	_ = vs.Set("proj", "LOG_LEVEL", "public")
	_ = vs.Set("proj", "DB_PASS", "secret")

	names, err := vs.ListByLevel("proj", "secret")
	if err != nil {
		t.Fatalf("ListByLevel: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 secret vars, got %d", len(names))
	}
}

func TestVisibilityEmptyProject(t *testing.T) {
	vs := newVisibilityStore(t)
	if err := vs.Set("", "VAR", "public"); err == nil {
		t.Error("expected error for empty project")
	}
	if _, err := vs.Get("", "VAR"); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestVisibilityEmptyVarName(t *testing.T) {
	vs := newVisibilityStore(t)
	if err := vs.Set("proj", "", "public"); err == nil {
		t.Error("expected error for empty var name")
	}
}

func TestVisibilityListByLevelInvalid(t *testing.T) {
	vs := newVisibilityStore(t)
	if _, err := vs.ListByLevel("proj", "unknown"); err == nil {
		t.Error("expected error for invalid level in ListByLevel")
	}
}
