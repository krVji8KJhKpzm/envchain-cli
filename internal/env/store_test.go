package env_test

import (
	"testing"

	"github.com/envchain-cli/internal/env"
	"github.com/envchain-cli/internal/keychain"
)

func newTestStore(t *testing.T, project string) *env.Store {
	t.Helper()
	kc := keychain.New()
	s, err := env.New(project, kc)
	if err != nil {
		t.Fatalf("env.New: %v", err)
	}
	return s
}

func TestStoreSetAndGet(t *testing.T) {
	s := newTestStore(t, "testproject")
	if err := s.Set("DB_URL", "postgres://localhost/test"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	val, err := s.Get("DB_URL")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "postgres://localhost/test" {
		t.Errorf("expected %q, got %q", "postgres://localhost/test", val)
	}
	_ = s.Delete("DB_URL")
}

func TestStoreDelete(t *testing.T) {
	s := newTestStore(t, "testproject")
	_ = s.Set("API_KEY", "secret")
	if err := s.Delete("API_KEY"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Get("API_KEY")
	if err == nil {
		t.Error("expected error after deletion, got nil")
	}
}

func TestStoreEmptyProject(t *testing.T) {
	kc := keychain.New()
	_, err := env.New("", kc)
	if err == nil {
		t.Error("expected error for empty project name")
	}
}

func TestStoreInvalidVarName(t *testing.T) {
	s := newTestStore(t, "testproject")
	if err := s.Set("", "value"); err == nil {
		t.Error("expected error for empty variable name")
	}
	if err := s.Set("BAD NAME", "value"); err == nil {
		t.Error("expected error for variable name with space")
	}
}

func TestStoreSeparateProjects(t *testing.T) {
	a := newTestStore(t, "proj-a")
	b := newTestStore(t, "proj-b")
	_ = a.Set("TOKEN", "aaa")
	_ = b.Set("TOKEN", "bbb")
	va, _ := a.Get("TOKEN")
	vb, _ := b.Get("TOKEN")
	if va == vb {
		t.Error("projects should have isolated variable namespaces")
	}
	_ = a.Delete("TOKEN")
	_ = b.Delete("TOKEN")
}
