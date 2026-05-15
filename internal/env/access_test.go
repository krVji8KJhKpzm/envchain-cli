package env

import (
	"testing"
	"time"
)

func newAccessStore() *AccessStore {
	return NewAccessStore(newMemKV())
}

func TestAccessRecord(t *testing.T) {
	s := newAccessStore()
	if err := s.Record("myapp", "API_KEY", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, err := s.List("myapp", "API_KEY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Actor != "alice" {
		t.Errorf("expected actor alice, got %s", entries[0].Actor)
	}
	if entries[0].VarName != "API_KEY" {
		t.Errorf("expected varName API_KEY, got %s", entries[0].VarName)
	}
}

func TestAccessListEmpty(t *testing.T) {
	s := newAccessStore()
	entries, err := s.List("myapp", "MISSING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty list, got %d entries", len(entries))
	}
}

func TestAccessMultipleRecords(t *testing.T) {
	s := newAccessStore()
	_ = s.Record("proj", "SECRET", "alice")
	time.Sleep(time.Millisecond)
	_ = s.Record("proj", "SECRET", "bob")
	entries, err := s.List("proj", "SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// newest first
	if entries[0].Actor != "bob" {
		t.Errorf("expected newest entry first (bob), got %s", entries[0].Actor)
	}
}

func TestAccessClear(t *testing.T) {
	s := newAccessStore()
	_ = s.Record("proj", "TOKEN", "alice")
	if err := s.Clear("proj", "TOKEN"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := s.List("proj", "TOKEN")
	if len(entries) != 0 {
		t.Errorf("expected empty after clear, got %d", len(entries))
	}
}

func TestAccessEmptyProject(t *testing.T) {
	s := newAccessStore()
	if err := s.Record("", "VAR", "alice"); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestAccessEmptyVarName(t *testing.T) {
	s := newAccessStore()
	if err := s.Record("proj", "", "alice"); err == nil {
		t.Error("expected error for empty varName")
	}
}
