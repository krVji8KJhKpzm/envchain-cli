package env

import (
	"testing"
)

func newReadonlyStore(t *testing.T) *ReadonlyStore {
	t.Helper()
	kc := newMockKeychain()
	s, err := NewReadonlyStore(kc, "myproject")
	if err != nil {
		t.Fatalf("NewReadonlyStore: %v", err)
	}
	return s
}

func TestReadonlySetAndIsReadonly(t *testing.T) {
	s := newReadonlyStore(t)
	if err := s.SetReadonly("API_KEY"); err != nil {
		t.Fatalf("SetReadonly: %v", err)
	}
	ok, err := s.IsReadonly("API_KEY")
	if err != nil {
		t.Fatalf("IsReadonly: %v", err)
	}
	if !ok {
		t.Error("expected API_KEY to be readonly")
	}
}

func TestReadonlyIsReadonlyNotFound(t *testing.T) {
	s := newReadonlyStore(t)
	ok, err := s.IsReadonly("MISSING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected MISSING to not be readonly")
	}
}

func TestReadonlyUnset(t *testing.T) {
	s := newReadonlyStore(t)
	_ = s.SetReadonly("DB_PASS")
	if err := s.Unset("DB_PASS"); err != nil {
		t.Fatalf("Unset: %v", err)
	}
	ok, _ := s.IsReadonly("DB_PASS")
	if ok {
		t.Error("expected DB_PASS to no longer be readonly after Unset")
	}
}

func TestReadonlyUnsetNotFound(t *testing.T) {
	s := newReadonlyStore(t)
	err := s.Unset("NEVER_SET")
	if err == nil {
		t.Error("expected error when unsetting a flag that was never set")
	}
}

func TestReadonlyListReadonly(t *testing.T) {
	s := newReadonlyStore(t)
	_ = s.SetReadonly("VAR_A")
	_ = s.SetReadonly("VAR_B")
	names, err := s.ListReadonly()
	if err != nil {
		t.Fatalf("ListReadonly: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 readonly vars, got %d", len(names))
	}
	seen := map[string]bool{}
	for _, n := range names {
		seen[n] = true
	}
	for _, want := range []string{"VAR_A", "VAR_B"} {
		if !seen[want] {
			t.Errorf("expected %q in list", want)
		}
	}
}

func TestReadonlyEmptyVarName(t *testing.T) {
	s := newReadonlyStore(t)
	if err := s.SetReadonly(""); err == nil {
		t.Error("expected error for empty variable name")
	}
}

func TestNewReadonlyStoreEmptyProject(t *testing.T) {
	kc := newMockKeychain()
	_, err := NewReadonlyStore(kc, "")
	if err == nil {
		t.Error("expected error for empty project name")
	}
}
