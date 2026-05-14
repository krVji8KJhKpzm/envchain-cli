package env

import (
	"testing"
)

func newGroupStore(t *testing.T) *GroupStore {
	t.Helper()
	kv := newTestStore(t)
	gs, err := NewGroupStore(kv)
	if err != nil {
		t.Fatalf("NewGroupStore: %v", err)
	}
	return gs
}

func TestAddAndListGroup(t *testing.T) {
	gs := newGroupStore(t)
	if err := gs.AddToGroup("proj", "backend", "DB_URL"); err != nil {
		t.Fatalf("AddToGroup: %v", err)
	}
	if err := gs.AddToGroup("proj", "backend", "DB_PASS"); err != nil {
		t.Fatalf("AddToGroup: %v", err)
	}
	members, err := gs.ListGroup("proj", "backend")
	if err != nil {
		t.Fatalf("ListGroup: %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestAddDuplicateToGroup(t *testing.T) {
	gs := newGroupStore(t)
	_ = gs.AddToGroup("proj", "g1", "VAR_A")
	_ = gs.AddToGroup("proj", "g1", "VAR_A")
	members, _ := gs.ListGroup("proj", "g1")
	if len(members) != 1 {
		t.Fatalf("expected 1 member after duplicate add, got %d", len(members))
	}
}

func TestRemoveFromGroup(t *testing.T) {
	gs := newGroupStore(t)
	_ = gs.AddToGroup("proj", "g1", "VAR_A")
	_ = gs.AddToGroup("proj", "g1", "VAR_B")
	if err := gs.RemoveFromGroup("proj", "g1", "VAR_A"); err != nil {
		t.Fatalf("RemoveFromGroup: %v", err)
	}
	members, _ := gs.ListGroup("proj", "g1")
	if len(members) != 1 || members[0] != "VAR_B" {
		t.Fatalf("expected [VAR_B], got %v", members)
	}
}

func TestListGroups(t *testing.T) {
	gs := newGroupStore(t)
	_ = gs.AddToGroup("proj", "alpha", "X")
	_ = gs.AddToGroup("proj", "beta", "Y")
	groups, err := gs.ListGroups("proj")
	if err != nil {
		t.Fatalf("ListGroups: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d: %v", len(groups), groups)
	}
	if groups[0] != "alpha" || groups[1] != "beta" {
		t.Fatalf("unexpected groups order: %v", groups)
	}
}

func TestListGroupEmpty(t *testing.T) {
	gs := newGroupStore(t)
	members, err := gs.ListGroup("proj", "nonexistent")
	if err != nil {
		t.Fatalf("expected no error for missing group, got %v", err)
	}
	if len(members) != 0 {
		t.Fatalf("expected empty members, got %v", members)
	}
}

func TestGroupEmptyProject(t *testing.T) {
	gs := newGroupStore(t)
	if err := gs.AddToGroup("", "g", "V"); err == nil {
		t.Fatal("expected error for empty project")
	}
	if _, err := gs.ListGroups(""); err == nil {
		t.Fatal("expected error for empty project in ListGroups")
	}
}

func TestNewGroupStoreNilKV(t *testing.T) {
	_, err := NewGroupStore(nil)
	if err == nil {
		t.Fatal("expected error for nil kv store")
	}
}
