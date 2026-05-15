package env

import (
	"testing"
)

func newNoteStore(t *testing.T) *NoteStore {
	t.Helper()
	kc := newTestKeychain()
	ns, err := NewNoteStore(kc)
	if err != nil {
		t.Fatalf("NewNoteStore: %v", err)
	}
	return ns
}

func TestNoteSetAndGet(t *testing.T) {
	ns := newNoteStore(t)
	if err := ns.Set("myproject", "API_KEY", "Used for external API calls"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := ns.Get("myproject", "API_KEY")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "Used for external API calls" {
		t.Errorf("got %q, want %q", got, "Used for external API calls")
	}
}

func TestNoteGetNotFound(t *testing.T) {
	ns := newNoteStore(t)
	got, err := ns.Get("myproject", "MISSING")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestNoteRemove(t *testing.T) {
	ns := newNoteStore(t)
	_ = ns.Set("proj", "SECRET", "some note")
	if err := ns.Remove("proj", "SECRET"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	got, _ := ns.Get("proj", "SECRET")
	if got != "" {
		t.Errorf("expected empty after remove, got %q", got)
	}
}

func TestNoteRemoveNotFound(t *testing.T) {
	ns := newNoteStore(t)
	if err := ns.Remove("proj", "NONEXISTENT"); err != nil {
		t.Errorf("Remove of missing key should not error: %v", err)
	}
}

func TestNoteEmptyProject(t *testing.T) {
	ns := newNoteStore(t)
	if err := ns.Set("", "VAR", "note"); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestNoteEmptyVarName(t *testing.T) {
	ns := newNoteStore(t)
	if err := ns.Set("proj", "", "note"); err == nil {
		t.Error("expected error for empty var name")
	}
}

func TestNoteNilKeychain(t *testing.T) {
	_, err := NewNoteStore(nil)
	if err == nil {
		t.Error("expected error for nil keychain")
	}
}

func TestNoteOverwrite(t *testing.T) {
	ns := newNoteStore(t)
	_ = ns.Set("proj", "KEY", "first note")
	_ = ns.Set("proj", "KEY", "updated note")
	got, err := ns.Get("proj", "KEY")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "updated note" {
		t.Errorf("got %q, want %q", got, "updated note")
	}
}
