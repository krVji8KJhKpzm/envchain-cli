package env

import (
	"testing"
)

func newSnapshotStore(t *testing.T) *SnapshotStore {
	t.Helper()
	base := newTestStore(t)
	return NewSnapshotStore(base)
}

func TestSnapshotTake(t *testing.T) {
	ss := newSnapshotStore(t)
	_ = ss.store.Set("proj", "FOO", "bar")
	_ = ss.store.Set("proj", "BAZ", "qux")

	snap, err := ss.Take("proj", "initial")
	if err != nil {
		t.Fatalf("Take: %v", err)
	}
	if snap.Project != "proj" {
		t.Errorf("project = %q, want %q", snap.Project, "proj")
	}
	if snap.Label != "initial" {
		t.Errorf("label = %q, want %q", snap.Label, "initial")
	}
	if snap.Vars["FOO"] != "bar" || snap.Vars["BAZ"] != "qux" {
		t.Errorf("vars mismatch: %v", snap.Vars)
	}
}

func TestSnapshotList(t *testing.T) {
	ss := newSnapshotStore(t)
	_ = ss.store.Set("proj", "KEY", "v1")
	_, _ = ss.Take("proj", "snap1")
	_ = ss.store.Set("proj", "KEY", "v2")
	_, _ = ss.Take("proj", "snap2")

	snaps := ss.List("proj")
	if len(snaps) != 2 {
		t.Fatalf("want 2 snapshots, got %d", len(snaps))
	}
	if snaps[0].Label != "snap1" {
		t.Errorf("first snap label = %q", snaps[0].Label)
	}
}

func TestSnapshotRestore(t *testing.T) {
	ss := newSnapshotStore(t)
	_ = ss.store.Set("proj", "KEY", "original")
	_, _ = ss.Take("proj", "before")
	_ = ss.store.Set("proj", "KEY", "changed")

	if err := ss.Restore("proj", 0); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	val, _ := ss.store.Get("proj", "KEY")
	if val != "original" {
		t.Errorf("after restore KEY = %q, want %q", val, "original")
	}
}

func TestSnapshotRestoreOutOfRange(t *testing.T) {
	ss := newSnapshotStore(t)
	if err := ss.Restore("proj", 0); err == nil {
		t.Error("expected error for out-of-range index")
	}
}

func TestSnapshotEmptyProject(t *testing.T) {
	ss := newSnapshotStore(t)
	_, err := ss.Take("", "")
	if err == nil {
		t.Error("expected error for empty project")
	}
}

func TestSnapshotListEmpty(t *testing.T) {
	ss := newSnapshotStore(t)
	snaps := ss.List("nonexistent")
	if len(snaps) != 0 {
		t.Errorf("want 0 snapshots, got %d", len(snaps))
	}
}
