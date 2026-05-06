package env

import (
	"errors"
	"testing"
)

// mockStore is an in-memory Store implementation for testing.
type mockStore struct {
	data map[string]string
}

func newMockStore() *mockStore {
	return &mockStore{data: make(map[string]string)}
}

func (m *mockStore) key(project, name string) string {
	return project + "/" + name
}

func (m *mockStore) Get(project, name string) (string, error) {
	v, ok := m.data[m.key(project, name)]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (m *mockStore) Set(project, name, value string) error {
	m.data[m.key(project, name)] = value
	return nil
}

func (m *mockStore) Delete(project, name string) error {
	k := m.key(project, name)
	if _, ok := m.data[k]; !ok {
		return errors.New("not found")
	}
	delete(m.data, k)
	return nil
}

func TestCopyVar(t *testing.T) {
	s := newMockStore()
	_ = s.Set("alpha", "DB_URL", "postgres://localhost/alpha")

	if err := CopyVar(s, s, "alpha", "beta", "DB_URL", ""); err != nil {
		t.Fatalf("CopyVar: %v", err)
	}

	v, err := s.Get("beta", "DB_URL")
	if err != nil || v != "postgres://localhost/alpha" {
		t.Errorf("expected copied value, got %q, err=%v", v, err)
	}

	// original should still exist
	if _, err := s.Get("alpha", "DB_URL"); err != nil {
		t.Error("original should still exist after copy")
	}
}

func TestCopyVarRename(t *testing.T) {
	s := newMockStore()
	_ = s.Set("alpha", "OLD_NAME", "value123")

	if err := CopyVar(s, s, "alpha", "alpha", "OLD_NAME", "NEW_NAME"); err != nil {
		t.Fatalf("CopyVar rename: %v", err)
	}

	v, err := s.Get("alpha", "NEW_NAME")
	if err != nil || v != "value123" {
		t.Errorf("expected renamed value, got %q, err=%v", v, err)
	}
}

func TestMoveVar(t *testing.T) {
	s := newMockStore()
	_ = s.Set("alpha", "SECRET", "topsecret")

	if err := MoveVar(s, "alpha", "beta", "SECRET", ""); err != nil {
		t.Fatalf("MoveVar: %v", err)
	}

	if _, err := s.Get("alpha", "SECRET"); err == nil {
		t.Error("original should be deleted after move")
	}

	v, err := s.Get("beta", "SECRET")
	if err != nil || v != "topsecret" {
		t.Errorf("expected moved value, got %q, err=%v", v, err)
	}
}

func TestCopyVarNotFound(t *testing.T) {
	s := newMockStore()
	if err := CopyVar(s, s, "alpha", "beta", "MISSING", ""); err == nil {
		t.Error("expected error when source variable not found")
	}
}
