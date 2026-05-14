package env

import (
	"testing"
)

func newAliasStore(t *testing.T) *AliasStore {
	t.Helper()
	kc := newMockKeychain()
	s, err := NewAliasStore(kc)
	if err != nil {
		t.Fatalf("NewAliasStore: %v", err)
	}
	return s
}

func TestAliasSetAndResolve(t *testing.T) {
	s := newAliasStore(t)
	if err := s.Set("mydb", "myproject", "DB_URL"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	proj, varName, err := s.Resolve("mydb")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if proj != "myproject" || varName != "DB_URL" {
		t.Errorf("got %q/%q, want myproject/DB_URL", proj, varName)
	}
}

func TestAliasResolveNotFound(t *testing.T) {
	s := newAliasStore(t)
	_, _, err := s.Resolve("ghost")
	if err == nil {
		t.Fatal("expected error for missing alias")
	}
}

func TestAliasRemove(t *testing.T) {
	s := newAliasStore(t)
	_ = s.Set("tok", "proj", "API_TOKEN")
	if err := s.Remove("tok"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, _, err := s.Resolve("tok")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestAliasRemoveNotFound(t *testing.T) {
	s := newAliasStore(t)
	if err := s.Remove("nope"); err == nil {
		t.Fatal("expected error removing non-existent alias")
	}
}

func TestAliasList(t *testing.T) {
	s := newAliasStore(t)
	_ = s.Set("a1", "p", "V1")
	_ = s.Set("a2", "p", "V2")
	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(list))
	}
}

func TestAliasInvalidName(t *testing.T) {
	s := newAliasStore(t)
	cases := []string{"", "123bad", "has space", "!nope"}
	for _, c := range cases {
		if err := s.Set(c, "proj", "VAR"); err == nil {
			t.Errorf("expected error for alias name %q", c)
		}
	}
}

func TestAliasNilKeychain(t *testing.T) {
	_, err := NewAliasStore(nil)
	if err == nil {
		t.Fatal("expected error for nil keychain")
	}
}
