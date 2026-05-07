package env

import (
	"testing"
	"time"
)

func newLockStore(t *testing.T) (*LockStore, *mockStore) {
	t.Helper()
	ms := newMockStore()
	ls, err := NewLockStore(ms, "testproject")
	if err != nil {
		t.Fatalf("NewLockStore: %v", err)
	}
	return ls, ms
}

func TestLockAndIsLocked(t *testing.T) {
	ls, _ := newLockStore(t)

	locked, err := ls.IsLocked("API_KEY")
	if err != nil || locked {
		t.Fatalf("expected not locked, got locked=%v err=%v", locked, err)
	}

	if err := ls.Lock("API_KEY", "alice"); err != nil {
		t.Fatalf("Lock: %v", err)
	}

	locked, err = ls.IsLocked("API_KEY")
	if err != nil {
		t.Fatalf("IsLocked: %v", err)
	}
	if !locked {
		t.Error("expected variable to be locked")
	}
}

func TestUnlock(t *testing.T) {
	ls, _ := newLockStore(t)

	_ = ls.Lock("DB_PASS", "bob")
	if err := ls.Unlock("DB_PASS"); err != nil {
		t.Fatalf("Unlock: %v", err)
	}

	locked, err := ls.IsLocked("DB_PASS")
	if err != nil || locked {
		t.Errorf("expected unlocked after Unlock, got locked=%v err=%v", locked, err)
	}
}

func TestGetLockEntry(t *testing.T) {
	ls, _ := newLockStore(t)
	before := time.Now().UTC().Truncate(time.Second)

	if err := ls.Lock("SECRET", "carol"); err != nil {
		t.Fatalf("Lock: %v", err)
	}

	entry, err := ls.GetLockEntry("SECRET")
	if err != nil {
		t.Fatalf("GetLockEntry: %v", err)
	}
	if entry.LockedBy != "carol" {
		t.Errorf("expected actor carol, got %s", entry.LockedBy)
	}
	if entry.VarName != "SECRET" {
		t.Errorf("expected varname SECRET, got %s", entry.VarName)
	}
	if entry.LockedAt.Before(before) {
		t.Errorf("lock timestamp %v is before test start %v", entry.LockedAt, before)
	}
}

func TestLockEmptyVarName(t *testing.T) {
	ls, _ := newLockStore(t)
	if err := ls.Lock("", "alice"); err == nil {
		t.Error("expected error for empty var name")
	}
}

func TestNewLockStoreEmptyProject(t *testing.T) {
	ms := newMockStore()
	_, err := NewLockStore(ms, "")
	if err == nil {
		t.Error("expected error for empty project")
	}
}

func TestUnlockNotLocked(t *testing.T) {
	ls, _ := newLockStore(t)
	// Unlocking a variable that was never locked should not panic
	// (may return error from underlying store, which is acceptable)
	_ = ls.Unlock("NONEXISTENT")
}
