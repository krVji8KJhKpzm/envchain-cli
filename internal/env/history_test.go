package env

import (
	"testing"
	"time"
)

func newHistoryStore(t *testing.T) *HistoryStore {
	t.Helper()
	s := newTestStore(t)
	return NewHistoryStore(s, 5)
}

func TestHistoryRecord(t *testing.T) {
	h := newHistoryStore(t)

	if err := h.Record("proj", "TOKEN", "abc123", "alice"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := h.List("proj", "TOKEN")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Value != "abc123" {
		t.Errorf("expected value abc123, got %s", entries[0].Value)
	}
	if entries[0].Actor != "alice" {
		t.Errorf("expected actor alice, got %s", entries[0].Actor)
	}
}

func TestHistoryNewestFirst(t *testing.T) {
	h := newHistoryStore(t)

	for _, v := range []string{"v1", "v2", "v3"} {
		if err := h.Record("proj", "KEY", v, ""); err != nil {
			t.Fatalf("Record %s: %v", v, err)
		}
	}

	entries, err := h.List("proj", "KEY")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if entries[0].Value != "v3" {
		t.Errorf("expected newest first, got %s", entries[0].Value)
	}
}

func TestHistoryMaxSize(t *testing.T) {
	s := newTestStore(t)
	h := NewHistoryStore(s, 3)

	for i := 0; i < 6; i++ {
		_ = h.Record("proj", "KEY", fmt.Sprintf("val%d", i), "")
	}

	entries, err := h.List("proj", "KEY")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected max 3 entries, got %d", len(entries))
	}
}

func TestHistoryListNotFound(t *testing.T) {
	h := newHistoryStore(t)
	_, err := h.List("proj", "MISSING")
	if err == nil {
		t.Error("expected error for missing history, got nil")
	}
}

func TestHistoryClear(t *testing.T) {
	h := newHistoryStore(t)
	_ = h.Record("proj", "KEY", "val", "bob")

	if err := h.Clear("proj", "KEY"); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	_, err := h.List("proj", "KEY")
	if err == nil {
		t.Error("expected error after clear, got nil")
	}
}

func TestHistoryTimestamp(t *testing.T) {
	h := newHistoryStore(t)
	before := time.Now().UTC().Add(-time.Second)
	_ = h.Record("proj", "KEY", "val", "")
	after := time.Now().UTC().Add(time.Second)

	entries, _ := h.List("proj", "KEY")
	if entries[0].ChangedAt.Before(before) || entries[0].ChangedAt.After(after) {
		t.Errorf("timestamp %v out of expected range", entries[0].ChangedAt)
	}
}
