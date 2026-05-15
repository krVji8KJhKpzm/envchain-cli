package env

import (
	"testing"
)

func newProtectStore(t *testing.T) (*ProtectStore, *mockStore) {
	t.Helper()
	ms := newMockStore()
	ps, err := NewProtectStore(ms, "testproject")
	if err != nil {
		t.Fatalf("NewProtectStore: %v", err)
	}
	return ps, ms
}

func TestProtectAndIsProtected(t *testing.T) {
	ps, _ := newProtectStore(t)
	if err := ps.Protect("API_KEY"); err != nil {
		t.Fatalf("Protect: %v", err)
	}
	ok, err := ps.IsProtected("API_KEY")
	if err != nil {
		t.Fatalf("IsProtected: %v", err)
	}
	if !ok {
		t.Error("expected API_KEY to be protected")
	}
}

func TestIsProtectedNotFound(t *testing.T) {
	ps, _ := newProtectStore(t)
	ok, err := ps.IsProtected("MISSING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected MISSING to not be protected")
	}
}

func TestUnprotect(t *testing.T) {
	ps, _ := newProtectStore(t)
	_ = ps.Protect("DB_PASS")
	if err := ps.Unprotect("DB_PASS"); err != nil {
		t.Fatalf("Unprotect: %v", err)
	}
	ok, _ := ps.IsProtected("DB_PASS")
	if ok {
		t.Error("expected DB_PASS to be unprotected after Unprotect")
	}
}

func TestListProtected(t *testing.T) {
	ps, _ := newProtectStore(t)
	_ = ps.Protect("TOKEN")
	_ = ps.Protect("SECRET")
	list, err := ps.ListProtected()
	if err != nil {
		t.Fatalf("ListProtected: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 protected vars, got %d", len(list))
	}
}

func TestProtectEmptyVarName(t *testing.T) {
	ps, _ := newProtectStore(t)
	if err := ps.Protect(""); err == nil {
		t.Error("expected error for empty variable name")
	}
}

func TestNewProtectStoreEmptyProject(t *testing.T) {
	ms := newMockStore()
	_, err := NewProtectStore(ms, "")
	if err == nil {
		t.Error("expected error for empty project name")
	}
}
