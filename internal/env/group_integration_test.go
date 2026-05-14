package env_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/env"
	"github.com/yourorg/envchain/internal/keychain"
)

func newIntegrationGroupStore(t *testing.T) (*env.GroupStore, *env.Store) {
	t.Helper()
	kc := keychain.NewMemory()
	store, err := env.New(kc)
	if err != nil {
		t.Fatalf("env.New: %v", err)
	}
	gs, err := env.NewGroupStore(store)
	if err != nil {
		t.Fatalf("NewGroupStore: %v", err)
	}
	return gs, store
}

func TestGroupIntegrationAddListRemove(t *testing.T) {
	gs, store := newIntegrationGroupStore(t)
	const proj = "myapp"

	_ = store.Set(proj, "API_KEY", "secret")
	_ = store.Set(proj, "DB_URL", "postgres://localhost")

	if err := gs.AddToGroup(proj, "infra", "API_KEY"); err != nil {
		t.Fatalf("AddToGroup API_KEY: %v", err)
	}
	if err := gs.AddToGroup(proj, "infra", "DB_URL"); err != nil {
		t.Fatalf("AddToGroup DB_URL: %v", err)
	}

	members, err := gs.ListGroup(proj, "infra")
	if err != nil {
		t.Fatalf("ListGroup: %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}

	if err := gs.RemoveFromGroup(proj, "infra", "API_KEY"); err != nil {
		t.Fatalf("RemoveFromGroup: %v", err)
	}
	members, _ = gs.ListGroup(proj, "infra")
	if len(members) != 1 || members[0] != "DB_URL" {
		t.Fatalf("expected [DB_URL] after remove, got %v", members)
	}
}

func TestGroupIntegrationMultipleGroups(t *testing.T) {
	gs, _ := newIntegrationGroupStore(t)
	const proj = "svc"

	_ = gs.AddToGroup(proj, "frontend", "REACT_APP_URL")
	_ = gs.AddToGroup(proj, "backend", "DB_PASS")
	_ = gs.AddToGroup(proj, "backend", "JWT_SECRET")

	groups, err := gs.ListGroups(proj)
	if err != nil {
		t.Fatalf("ListGroups: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d: %v", len(groups), groups)
	}

	backend, _ := gs.ListGroup(proj, "backend")
	if len(backend) != 2 {
		t.Fatalf("expected 2 backend members, got %d", len(backend))
	}
}
