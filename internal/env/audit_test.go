package env

import (
	"strings"
	"testing"
	"time"
)

func TestAuditRecord(t *testing.T) {
	log := NewAuditLog()
	log.Record("myproject", "API_KEY", "set", "alice")

	if log.Len() != 1 {
		t.Fatalf("expected 1 event, got %d", log.Len())
	}

	events := log.Events()
	e := events[0]
	if e.Project != "myproject" {
		t.Errorf("expected project 'myproject', got %q", e.Project)
	}
	if e.VarName != "API_KEY" {
		t.Errorf("expected var 'API_KEY', got %q", e.VarName)
	}
	if e.Action != "set" {
		t.Errorf("expected action 'set', got %q", e.Action)
	}
	if e.Actor != "alice" {
		t.Errorf("expected actor 'alice', got %q", e.Actor)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestAuditEventString(t *testing.T) {
	e := AuditEvent{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Project:   "proj",
		VarName:   "SECRET",
		Action:    "delete",
		Actor:     "bob",
	}
	s := e.String()
	if !strings.Contains(s, "delete") {
		t.Errorf("expected 'delete' in string, got %q", s)
	}
	if !strings.Contains(s, "proj/SECRET") {
		t.Errorf("expected 'proj/SECRET' in string, got %q", s)
	}
	if !strings.Contains(s, "bob") {
		t.Errorf("expected 'bob' in string, got %q", s)
	}
}

func TestAuditFilterByProject(t *testing.T) {
	log := NewAuditLog()
	log.Record("alpha", "KEY1", "set", "alice")
	log.Record("beta", "KEY2", "set", "bob")
	log.Record("alpha", "KEY3", "delete", "alice")

	result := log.FilterByProject("alpha")
	if len(result) != 2 {
		t.Fatalf("expected 2 events for 'alpha', got %d", len(result))
	}
	for _, e := range result {
		if e.Project != "alpha" {
			t.Errorf("unexpected project %q in filtered results", e.Project)
		}
	}
}

func TestAuditFilterByAction(t *testing.T) {
	log := NewAuditLog()
	log.Record("proj", "K1", "set", "alice")
	log.Record("proj", "K2", "rotate", "alice")
	log.Record("proj", "K3", "set", "alice")

	result := log.FilterByAction("set")
	if len(result) != 2 {
		t.Fatalf("expected 2 'set' events, got %d", len(result))
	}
}

func TestAuditEventsReturnsCopy(t *testing.T) {
	log := NewAuditLog()
	log.Record("proj", "KEY", "set", "alice")

	events := log.Events()
	events[0].Actor = "tampered"

	original := log.Events()
	if original[0].Actor == "tampered" {
		t.Error("Events() should return a copy, not a reference")
	}
}

func TestAuditEmptyLog(t *testing.T) {
	log := NewAuditLog()
	if log.Len() != 0 {
		t.Errorf("expected empty log, got %d events", log.Len())
	}
	if events := log.Events(); len(events) != 0 {
		t.Errorf("expected empty events slice")
	}
}
