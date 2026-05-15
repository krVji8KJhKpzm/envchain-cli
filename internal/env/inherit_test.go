package env

import (
	"testing"
)

func newInheritStore(t *testing.T) (*InheritStore, Store) {
	t.Helper()
	s := newTestStore(t)
	return NewInheritStore(s), s
}

func TestInheritSetAndGetParent(t *testing.T) {
	h, _ := newInheritStore(t)
	if err := h.SetParent("child", "parent"); err != nil {
		t.Fatalf("SetParent: %v", err)
	}
	got, err := h.GetParent("child")
	if err != nil {
		t.Fatalf("GetParent: %v", err)
	}
	if got != "parent" {
		t.Errorf("expected parent, got %q", got)
	}
}

func TestInheritGetParentNotSet(t *testing.T) {
	h, _ := newInheritStore(t)
	got, err := h.GetParent("orphan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty parent, got %q", got)
	}
}

func TestInheritRemoveParent(t *testing.T) {
	h, _ := newInheritStore(t)
	_ = h.SetParent("child", "parent")
	if err := h.RemoveParent("child"); err != nil {
		t.Fatalf("RemoveParent: %v", err)
	}
	got, _ := h.GetParent("child")
	if got != "" {
		t.Errorf("expected empty after removal, got %q", got)
	}
}

func TestInheritResolveChildOverridesParent(t *testing.T) {
	h, s := newInheritStore(t)
	_ = s.Set("base", "API_KEY", "base-key")
	_ = s.Set("child", "API_KEY", "child-key")
	_ = h.SetParent("child", "base")
	v, err := h.Resolve("child", "API_KEY")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if v != "child-key" {
		t.Errorf("expected child-key, got %q", v)
	}
}

func TestInheritResolveFallsBackToParent(t *testing.T) {
	h, s := newInheritStore(t)
	_ = s.Set("base", "DB_URL", "postgres://localhost")
	_ = h.SetParent("child", "base")
	v, err := h.Resolve("child", "DB_URL")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if v != "postgres://localhost" {
		t.Errorf("expected postgres://localhost, got %q", v)
	}
}

func TestInheritResolveNotFound(t *testing.T) {
	h, _ := newInheritStore(t)
	_, err := h.Resolve("child", "MISSING")
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestInheritSetParentSelf(t *testing.T) {
	h, _ := newInheritStore(t)
	if err := h.SetParent("proj", "proj"); err == nil {
		t.Fatal("expected error when child == parent")
	}
}

func TestInheritSetParentEmptyChild(t *testing.T) {
	h, _ := newInheritStore(t)
	if err := h.SetParent("", "parent"); err == nil {
		t.Fatal("expected error for empty child")
	}
}
