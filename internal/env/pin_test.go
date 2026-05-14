package env

import (
	"testing"
)

func newPinStore(t *testing.T) (*PinStore, *mockStore) {
	t.Helper()
	ms := newMockStore()
	ps, err := NewPinStore(ms, "myproject")
	if err != nil {
		t.Fatalf("NewPinStore: %v", err)
	}
	return ps, ms
}

func TestPinAndIsPinned(t *testing.T) {
	ps, _ := newPinStore(t)

	if err := ps.Pin("API_KEY"); err != nil {
		t.Fatalf("Pin: %v", err)
	}

	pinned, err := ps.IsPinned("API_KEY")
	if err != nil {
		t.Fatalf("IsPinned: %v", err)
	}
	if !pinned {
		t.Error("expected API_KEY to be pinned")
	}
}

func TestIsPinnedNotFound(t *testing.T) {
	ps, _ := newPinStore(t)

	pinned, err := ps.IsPinned("MISSING")
	if err != nil {
		t.Fatalf("IsPinned: %v", err)
	}
	if pinned {
		t.Error("expected MISSING to not be pinned")
	}
}

func TestUnpin(t *testing.T) {
	ps, _ := newPinStore(t)

	_ = ps.Pin("DB_PASS")
	if err := ps.Unpin("DB_PASS"); err != nil {
		t.Fatalf("Unpin: %v", err)
	}

	pinned, _ := ps.IsPinned("DB_PASS")
	if pinned {
		t.Error("expected DB_PASS to be unpinned")
	}
}

func TestListPinned(t *testing.T) {
	ps, _ := newPinStore(t)

	_ = ps.Pin("A")
	_ = ps.Pin("B")
	_ = ps.Pin("C")

	list, err := ps.ListPinned()
	if err != nil {
		t.Fatalf("ListPinned: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 pinned vars, got %d", len(list))
	}
}

func TestPinEmptyVarName(t *testing.T) {
	ps, _ := newPinStore(t)

	if err := ps.Pin(""); err == nil {
		t.Error("expected error for empty var name")
	}
}

func TestNewPinStoreEmptyProject(t *testing.T) {
	ms := newMockStore()
	_, err := NewPinStore(ms, "")
	if err == nil {
		t.Error("expected error for empty project")
	}
}
