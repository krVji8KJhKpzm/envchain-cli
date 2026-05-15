package env

import (
	"testing"
)

func newCategoryStore(t *testing.T) *CategoryStore {
	t.Helper()
	kc := newTestKeychain()
	return NewCategoryStore(kc)
}

func TestCategorySetAndGet(t *testing.T) {
	s := newCategoryStore(t)
	if err := s.Set("proj", "API_KEY", "secrets"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	cat, err := s.Get("proj", "API_KEY")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if cat != "secrets" {
		t.Errorf("expected %q, got %q", "secrets", cat)
	}
}

func TestCategoryGetNotFound(t *testing.T) {
	s := newCategoryStore(t)
	cat, err := s.Get("proj", "MISSING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat != "" {
		t.Errorf("expected empty string, got %q", cat)
	}
}

func TestCategoryRemove(t *testing.T) {
	s := newCategoryStore(t)
	_ = s.Set("proj", "DB_PASS", "database")
	if err := s.Remove("proj", "DB_PASS"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	cat, _ := s.Get("proj", "DB_PASS")
	if cat != "" {
		t.Errorf("expected empty after remove, got %q", cat)
	}
}

func TestCategoryInvalidCategory(t *testing.T) {
	s := newCategoryStore(t)
	err := s.Set("proj", "VAR", "bad category!")
	if err == nil {
		t.Error("expected error for invalid category, got nil")
	}
}

func TestCategoryEmptyProject(t *testing.T) {
	s := newCategoryStore(t)
	if err := s.Set("", "VAR", "ok"); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestCategoryListByCategory(t *testing.T) {
	s := newCategoryStore(t)
	_ = s.Set("proj", "API_KEY", "secrets")
	_ = s.Set("proj", "DB_PASS", "secrets")
	_ = s.Set("proj", "PORT", "config")

	allVars := []string{"API_KEY", "DB_PASS", "PORT", "UNCAT"}
	result, err := s.ListByCategory("proj", "secrets", allVars)
	if err != nil {
		t.Fatalf("ListByCategory: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0] != "API_KEY" || result[1] != "DB_PASS" {
		t.Errorf("unexpected results: %v", result)
	}
}
