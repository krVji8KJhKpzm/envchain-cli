package env

import (
	"testing"
)

func newLabelStore(t *testing.T) (*LabelStore, Store) {
	t.Helper()
	s := newTestStore(t)
	return NewLabelStore(s), s
}

func TestLabelSetAndGet(t *testing.T) {
	ls, _ := newLabelStore(t)
	if err := ls.Set("proj", "API_KEY", "owner", "alice"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	val, err := ls.Get("proj", "API_KEY", "owner")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "alice" {
		t.Errorf("expected alice, got %q", val)
	}
}

func TestLabelGetNotFound(t *testing.T) {
	ls, _ := newLabelStore(t)
	_, err := ls.Get("proj", "API_KEY", "missing")
	if err == nil {
		t.Fatal("expected error for missing label")
	}
}

func TestLabelRemove(t *testing.T) {
	ls, _ := newLabelStore(t)
	_ = ls.Set("proj", "API_KEY", "env", "production")
	if err := ls.Remove("proj", "API_KEY", "env"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, err := ls.Get("proj", "API_KEY", "env")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestLabelList(t *testing.T) {
	ls, _ := newLabelStore(t)
	_ = ls.Set("proj", "DB_PASS", "owner", "bob")
	_ = ls.Set("proj", "DB_PASS", "team", "backend")
	_ = ls.Set("proj", "OTHER_VAR", "owner", "carol")

	labels, err := ls.List("proj", "DB_PASS")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
	if labels["owner"] != "bob" {
		t.Errorf("expected owner=bob, got %q", labels["owner"])
	}
	if labels["team"] != "backend" {
		t.Errorf("expected team=backend, got %q", labels["team"])
	}
}

func TestLabelEmptyArgs(t *testing.T) {
	ls, _ := newLabelStore(t)
	if err := ls.Set("", "VAR", "k", "v"); err == nil {
		t.Error("expected error for empty project")
	}
	if err := ls.Set("proj", "", "k", "v"); err == nil {
		t.Error("expected error for empty varName")
	}
	if err := ls.Set("proj", "VAR", "", "v"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestLabelKeyNoWhitespace(t *testing.T) {
	ls, _ := newLabelStore(t)
	if err := ls.Set("proj", "VAR", "bad key", "v"); err == nil {
		t.Error("expected error for whitespace in key")
	}
}
