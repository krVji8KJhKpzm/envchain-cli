package env

import (
	"testing"
)

func newNamespaceStore(t *testing.T) *NamespaceStore {
	t.Helper()
	s := newTestStore(t)
	ns, err := NewNamespaceStore(s)
	if err != nil {
		t.Fatalf("NewNamespaceStore: %v", err)
	}
	return ns
}

func TestProjectKey(t *testing.T) {
	ns := newNamespaceStore(t)
	key, err := ns.ProjectKey("acme", "payments")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "acme/payments" {
		t.Errorf("expected acme/payments, got %s", key)
	}
}

func TestProjectKeyEmptyNamespace(t *testing.T) {
	ns := newNamespaceStore(t)
	_, err := ns.ProjectKey("", "payments")
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestProjectKeyEmptyProject(t *testing.T) {
	ns := newNamespaceStore(t)
	_, err := ns.ProjectKey("acme", "")
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestProjectKeyInvalidNamespace(t *testing.T) {
	ns := newNamespaceStore(t)
	_, err := ns.ProjectKey("acme corp", "payments")
	if err == nil {
		t.Fatal("expected error for namespace with space")
	}
}

func TestListProjects(t *testing.T) {
	ns := newNamespaceStore(t)
	all := []string{"acme/payments", "acme/auth", "other/service", "acme/infra"}
	projects, err := ns.ListProjects("acme", all)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 3 {
		t.Fatalf("expected 3 projects, got %d", len(projects))
	}
	for _, p := range projects {
		if p == "payments" || p == "auth" || p == "infra" {
			continue
		}
		t.Errorf("unexpected project: %s", p)
	}
}

func TestListProjectsEmpty(t *testing.T) {
	ns := newNamespaceStore(t)
	projects, err := ns.ListProjects("acme", []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 0 {
		t.Errorf("expected 0 projects, got %d", len(projects))
	}
}

func TestParseProjectKey(t *testing.T) {
	ns, proj, err := ParseProjectKey("acme/payments")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ns != "acme" || proj != "payments" {
		t.Errorf("expected acme/payments, got %s/%s", ns, proj)
	}
}

func TestParseProjectKeyNoSlash(t *testing.T) {
	_, _, err := ParseProjectKey("noslash")
	if err == nil {
		t.Fatal("expected error for key without slash")
	}
}

func TestNewNamespaceStoreNilStore(t *testing.T) {
	_, err := NewNamespaceStore(nil)
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}
