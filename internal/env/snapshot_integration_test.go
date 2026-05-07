package env_test

import (
	"testing"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func newIntegrationSnapshotStore(t *testing.T) *env.SnapshotStore {
	t.Helper()
	kc, err := keychain.New("envchain-test-snapshot-" + t.Name())
	if err != nil {
		t.Skipf("keychain unavailable: %v", err)
	}
	t.Cleanup(func() { _ = kc.DeleteAll() })
	return env.NewSnapshotStore(env.New(kc))
}

func TestSnapshotIntegrationTakeAndList(t *testing.T) {
	ss := newIntegrationSnapshotStore(t)
	_ = ss.Store().Set("integ", "ALPHA", "one")
	_ = ss.Store().Set("integ", "BETA", "two")

	_, err := ss.Take("integ", "v1")
	if err != nil {
		t.Fatalf("Take: %v", err)
	}
	snaps := ss.List("integ")
	if len(snaps) != 1 {
		t.Fatalf("want 1 snapshot, got %d", len(snaps))
	}
	if snaps[0].Vars["ALPHA"] != "one" {
		t.Errorf("ALPHA = %q, want %q", snaps[0].Vars["ALPHA"], "one")
	}
}

func TestSnapshotIntegrationRestore(t *testing.T) {
	ss := newIntegrationSnapshotStore(t)
	_ = ss.Store().Set("integ", "KEY", "before")
	_, _ = ss.Take("integ", "checkpoint")
	_ = ss.Store().Set("integ", "KEY", "after")

	if err := ss.Restore("integ", 0); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	val, err := ss.Store().Get("integ", "KEY")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "before" {
		t.Errorf("KEY = %q, want %q", val, "before")
	}
}

func TestSnapshotIntegrationMultipleSnaps(t *testing.T) {
	ss := newIntegrationSnapshotStore(t)
	for i, v := range []string{"a", "b", "c"} {
		_ = ss.Store().Set("integ", "X", v)
		_, err := ss.Take("integ", v)
		if err != nil {
			t.Fatalf("Take %d: %v", i, err)
		}
	}
	snaps := ss.List("integ")
	if len(snaps) != 3 {
		t.Fatalf("want 3 snapshots, got %d", len(snaps))
	}
	if snaps[2].Label != "c" {
		t.Errorf("last snap label = %q, want %q", snaps[2].Label, "c")
	}
}
