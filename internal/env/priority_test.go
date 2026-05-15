package env

import (
	"testing"
)

func newPriorityStore(t *testing.T) *PriorityStore {
	t.Helper()
	kc := newTestKeychain(t)
	return NewPriorityStore(kc)
}

func TestPrioritySetAndGet(t *testing.T) {
	s := newPriorityStore(t)
	if err := s.Set("myproject", "API_KEY", 80); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := s.Get("myproject", "API_KEY")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != 80 {
		t.Errorf("expected 80, got %d", got)
	}
}

func TestPriorityDefaultWhenNotSet(t *testing.T) {
	s := newPriorityStore(t)
	got, err := s.Get("myproject", "MISSING")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != defaultPriority {
		t.Errorf("expected default %d, got %d", defaultPriority, got)
	}
}

func TestPriorityInvalidRange(t *testing.T) {
	s := newPriorityStore(t)
	if err := s.Set("proj", "VAR", 0); err == nil {
		t.Error("expected error for priority 0")
	}
	if err := s.Set("proj", "VAR", 101); err == nil {
		t.Error("expected error for priority 101")
	}
}

func TestPriorityEmptyProject(t *testing.T) {
	s := newPriorityStore(t)
	if err := s.Set("", "VAR", 50); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestPriorityRemove(t *testing.T) {
	s := newPriorityStore(t)
	if err := s.Set("proj", "TOKEN", 75); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.Remove("proj", "TOKEN"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	got, err := s.Get("proj", "TOKEN")
	if err != nil {
		t.Fatalf("Get after remove: %v", err)
	}
	if got != defaultPriority {
		t.Errorf("expected default after remove, got %d", got)
	}
}

func TestPriorityRemoveNotFound(t *testing.T) {
	s := newPriorityStore(t)
	if err := s.Remove("proj", "NONEXISTENT"); err != nil {
		t.Errorf("Remove of nonexistent should not error, got: %v", err)
	}
}

func TestPriorityCompare(t *testing.T) {
	s := newPriorityStore(t)
	_ = s.Set("proj", "HIGH", 90)
	_ = s.Set("proj", "LOW", 10)

	if got := s.Compare("proj", "HIGH", "LOW"); got != 1 {
		t.Errorf("HIGH vs LOW: expected 1, got %d", got)
	}
	if got := s.Compare("proj", "LOW", "HIGH"); got != -1 {
		t.Errorf("LOW vs HIGH: expected -1, got %d", got)
	}
	if got := s.Compare("proj", "HIGH", "HIGH"); got != 0 {
		t.Errorf("HIGH vs HIGH: expected 0, got %d", got)
	}
}
