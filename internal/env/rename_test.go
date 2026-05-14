package env

import (
	"errors"
	"testing"
)

func newRenameStore(t *testing.T) Store {
	t.Helper()
	s, err := New(newTestKeychain(), "rename-test")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestRenameVar(t *testing.T) {
	store := newRenameStore(t)
	if err := store.Set("proj", "OLD_KEY", "secret"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	r, _ := NewRenamer(store)
	if err := r.RenameVar("proj", "OLD_KEY", "NEW_KEY"); err != nil {
		t.Fatalf("RenameVar: %v", err)
	}
	if _, err := store.Get("proj", "OLD_KEY"); err == nil {
		t.Error("expected OLD_KEY to be deleted")
	}
	v, err := store.Get("proj", "NEW_KEY")
	if err != nil {
		t.Fatalf("Get NEW_KEY: %v", err)
	}
	if v != "secret" {
		t.Errorf("expected %q, got %q", "secret", v)
	}
}

func TestRenameVarSourceNotFound(t *testing.T) {
	store := newRenameStore(t)
	r, _ := NewRenamer(store)
	err := r.RenameVar("proj", "MISSING", "NEW_KEY")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRenameVarDestinationExists(t *testing.T) {
	store := newRenameStore(t)
	_ = store.Set("proj", "KEY_A", "val_a")
	_ = store.Set("proj", "KEY_B", "val_b")
	r, _ := NewRenamer(store)
	err := r.RenameVar("proj", "KEY_A", "KEY_B")
	if err == nil {
		t.Fatal("expected error when destination already exists")
	}
}

func TestRenameVarSameName(t *testing.T) {
	store := newRenameStore(t)
	_ = store.Set("proj", "KEY", "v")
	r, _ := NewRenamer(store)
	if err := r.RenameVar("proj", "KEY", "KEY"); err == nil {
		t.Fatal("expected error for identical names")
	}
}

func TestRenameVarInvalidName(t *testing.T) {
	store := newRenameStore(t)
	_ = store.Set("proj", "VALID", "v")
	r, _ := NewRenamer(store)
	if err := r.RenameVar("proj", "VALID", "123INVALID"); err == nil {
		t.Fatal("expected error for invalid destination name")
	}
}

func TestRenameVarEmptyProject(t *testing.T) {
	store := newRenameStore(t)
	r, _ := NewRenamer(store)
	if err := r.RenameVar("", "A", "B"); err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestNewRenamerNilStore(t *testing.T) {
	_, err := NewRenamer(nil)
	if !errors.Is(err, err) || err == nil {
		t.Fatal("expected error for nil store")
	}
}
