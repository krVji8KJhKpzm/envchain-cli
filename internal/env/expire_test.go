package env

import (
	"testing"
	"time"
)

func newExpiryStore(t *testing.T) *ExpiryStore {
	t.Helper()
	kc := newMockKeychain()
	s, err := NewExpiryStore(kc, "testproject")
	if err != nil {
		t.Fatalf("NewExpiryStore: %v", err)
	}
	return s
}

func TestExpirySetAndGet(t *testing.T) {
	s := newExpiryStore(t)
	at := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	if err := s.SetExpiry("API_KEY", at); err != nil {
		t.Fatalf("SetExpiry: %v", err)
	}
	entry, err := s.GetExpiry("API_KEY")
	if err != nil {
		t.Fatalf("GetExpiry: %v", err)
	}
	if !entry.ExpiresAt.Equal(at.UTC()) {
		t.Errorf("got %v, want %v", entry.ExpiresAt, at.UTC())
	}
	if entry.IsExpired() {
		t.Error("expected entry to not be expired")
	}
}

func TestExpiryGetNotFound(t *testing.T) {
	s := newExpiryStore(t)
	_, err := s.GetExpiry("MISSING")
	if err == nil {
		t.Fatal("expected error for missing expiry")
	}
}

func TestExpiryIsExpired(t *testing.T) {
	s := newExpiryStore(t)
	past := time.Now().Add(-1 * time.Hour)
	if err := s.SetExpiry("OLD_TOKEN", past); err != nil {
		t.Fatalf("SetExpiry: %v", err)
	}
	entry, err := s.GetExpiry("OLD_TOKEN")
	if err != nil {
		t.Fatalf("GetExpiry: %v", err)
	}
	if !entry.IsExpired() {
		t.Error("expected entry to be expired")
	}
}

func TestExpiryRemove(t *testing.T) {
	s := newExpiryStore(t)
	at := time.Now().Add(time.Hour)
	_ = s.SetExpiry("DB_PASS", at)
	if err := s.RemoveExpiry("DB_PASS"); err != nil {
		t.Fatalf("RemoveExpiry: %v", err)
	}
	_, err := s.GetExpiry("DB_PASS")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestExpiryListExpired(t *testing.T) {
	s := newExpiryStore(t)
	_ = s.SetExpiry("FRESH", time.Now().Add(time.Hour))
	_ = s.SetExpiry("STALE", time.Now().Add(-time.Hour))
	_ = s.SetExpiry("ANCIENT", time.Now().Add(-48*time.Hour))

	expired, err := s.ListExpired([]string{"FRESH", "STALE", "ANCIENT", "NO_EXPIRY"})
	if err != nil {
		t.Fatalf("ListExpired: %v", err)
	}
	if len(expired) != 2 {
		t.Errorf("got %d expired entries, want 2", len(expired))
	}
}

func TestExpiryEmptyProject(t *testing.T) {
	kc := newMockKeychain()
	_, err := NewExpiryStore(kc, "")
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestExpiryEntryString(t *testing.T) {
	future := ExpiryEntry{VarName: "TOKEN", ExpiresAt: time.Now().Add(time.Hour)}
	if s := future.String(); s == "" {
		t.Error("expected non-empty string for future entry")
	}
	past := ExpiryEntry{VarName: "TOKEN", ExpiresAt: time.Now().Add(-time.Hour)}
	if s := past.String(); s == "" {
		t.Error("expected non-empty string for past entry")
	}
}
