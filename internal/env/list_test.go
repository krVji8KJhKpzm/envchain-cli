package env

import (
	"testing"
)

func TestListVars(t *testing.T) {
	s := newTestStore(t)

	if err := s.Set("myapp", "DB_URL", "postgres://localhost"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.Set("myapp", "API_KEY", "secret"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.Set("myapp", "PORT", "8080"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	names, err := s.ListVars("myapp")
	if err != nil {
		t.Fatalf("ListVars: %v", err)
	}

	// Expect sorted order.
	expected := []string{"API_KEY", "DB_URL", "PORT"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d vars, got %d", len(expected), len(names))
	}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("names[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

func TestListVarsEmpty(t *testing.T) {
	s := newTestStore(t)

	names, err := s.ListVars("emptyproject")
	if err != nil {
		t.Fatalf("ListVars on empty project: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected 0 vars, got %d", len(names))
	}
}

func TestListVarsInvalidProject(t *testing.T) {
	s := newTestStore(t)

	_, err := s.ListVars("")
	if err == nil {
		t.Error("expected error for empty project name, got nil")
	}
}

func TestListAll(t *testing.T) {
	s := newTestStore(t)

	_ = s.Set("alpha", "FOO", "1")
	_ = s.Set("alpha", "BAR", "2")
	_ = s.Set("beta", "BAZ", "3")

	infos, err := s.ListAll([]string{"alpha", "beta"})
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}

	if len(infos) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(infos))
	}

	// First two belong to alpha (sorted), last to beta.
	if infos[0].Project != "alpha" || infos[1].Project != "alpha" {
		t.Errorf("expected first two entries from alpha")
	}
	if infos[2].Project != "beta" || infos[2].Name != "BAZ" {
		t.Errorf("unexpected last entry: %v", infos[2])
	}
}

func TestVarInfoString(t *testing.T) {
	v := VarInfo{Project: "myapp", Name: "DB_URL"}
	if got := v.String(); got != "myapp/DB_URL" {
		t.Errorf("String() = %q, want %q", got, "myapp/DB_URL")
	}
}
