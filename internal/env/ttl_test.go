package env

import (
	"strings"
	"testing"
	"time"
)

func TestTTLSetAndGet(t *testing.T) {
	s := NewTTLStore()
	s.Set("myproject", "API_KEY", 10*time.Minute)

	entry, ok := s.Get("myproject", "API_KEY")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if entry.Project != "myproject" || entry.VarName != "API_KEY" {
		t.Errorf("unexpected entry: %+v", entry)
	}
	if entry.IsExpired() {
		t.Error("entry should not be expired yet")
	}
}

func TestTTLGetNotFound(t *testing.T) {
	s := NewTTLStore()
	_, ok := s.Get("proj", "MISSING")
	if ok {
		t.Error("expected not found")
	}
}

func TestTTLIsExpired(t *testing.T) {
	s := NewTTLStore()
	s.Set("proj", "OLD_KEY", -1*time.Second)

	entry, ok := s.Get("proj", "OLD_KEY")
	if !ok {
		t.Fatal("expected entry")
	}
	if !entry.IsExpired() {
		t.Error("entry should be expired")
	}
}

func TestTTLExpiredList(t *testing.T) {
	s := NewTTLStore()
	s.Set("proj", "FRESH", 10*time.Minute)
	s.Set("proj", "STALE", -1*time.Second)

	expired := s.Expired()
	if len(expired) != 1 {
		t.Fatalf("expected 1 expired entry, got %d", len(expired))
	}
	if expired[0].VarName != "STALE" {
		t.Errorf("unexpected expired var: %s", expired[0].VarName)
	}
}

func TestTTLRemove(t *testing.T) {
	s := NewTTLStore()
	s.Set("proj", "API_KEY", 10*time.Minute)
	s.Remove("proj", "API_KEY")

	_, ok := s.Get("proj", "API_KEY")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestTTLUpdateExisting(t *testing.T) {
	s := NewTTLStore()
	s.Set("proj", "KEY", -1*time.Second)
	s.Set("proj", "KEY", 10*time.Minute) // refresh

	entry, _ := s.Get("proj", "KEY")
	if entry.IsExpired() {
		t.Error("entry should not be expired after refresh")
	}
	if len(s.entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(s.entries))
	}
}

func TestTTLEntryString(t *testing.T) {
	s := NewTTLStore()
	s.Set("proj", "KEY", 5*time.Minute)
	entry, _ := s.Get("proj", "KEY")

	str := entry.String()
	if !strings.Contains(str, "proj/KEY") {
		t.Errorf("expected project/var in string, got: %s", str)
	}
	if !strings.Contains(str, "expires in") {
		t.Errorf("expected 'expires in' in string, got: %s", str)
	}
}

func TestTTLEntryStringExpired(t *testing.T) {
	s := NewTTLStore()
	s.Set("proj", "OLD", -1*time.Second)
	entry, _ := s.Get("proj", "OLD")

	str := entry.String()
	if !strings.Contains(str, "EXPIRED") {
		t.Errorf("expected EXPIRED in string, got: %s", str)
	}
}
