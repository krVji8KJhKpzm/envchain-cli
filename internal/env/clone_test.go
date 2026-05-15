package env

import (
	"testing"
)

func newCloneStore() *mockStore {
	return newMockStore()
}

func TestCloneProject(t *testing.T) {
	s := newCloneStore()
	_ = s.Set("alpha", "DB_URL", "postgres://localhost")
	_ = s.Set("alpha", "API_KEY", "secret")

	c := NewCloner(s)
	res, err := c.CloneProject("alpha", "beta", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(res.Skipped))
	}

	val, err := s.Get("beta", "DB_URL")
	if err != nil || val != "postgres://localhost" {
		t.Errorf("expected DB_URL to be cloned, got %q %v", val, err)
	}
}

func TestCloneProjectSkipsExisting(t *testing.T) {
	s := newCloneStore()
	_ = s.Set("alpha", "DB_URL", "postgres://localhost")
	_ = s.Set("alpha", "API_KEY", "secret")
	_ = s.Set("beta", "DB_URL", "mysql://localhost")

	c := NewCloner(s)
	res, err := c.CloneProject("alpha", "beta", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 1 {
		t.Errorf("expected 1 copied, got %d", len(res.Copied))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}

	val, _ := s.Get("beta", "DB_URL")
	if val != "mysql://localhost" {
		t.Errorf("expected original value preserved, got %q", val)
	}
}

func TestCloneProjectOverwrite(t *testing.T) {
	s := newCloneStore()
	_ = s.Set("alpha", "DB_URL", "postgres://localhost")
	_ = s.Set("beta", "DB_URL", "mysql://localhost")

	c := NewCloner(s)
	res, err := c.CloneProject("alpha", "beta", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 1 {
		t.Errorf("expected 1 copied, got %d", len(res.Copied))
	}

	val, _ := s.Get("beta", "DB_URL")
	if val != "postgres://localhost" {
		t.Errorf("expected overwritten value, got %q", val)
	}
}

func TestCloneProjectSameSourceDest(t *testing.T) {
	s := newCloneStore()
	c := NewCloner(s)
	_, err := c.CloneProject("alpha", "alpha", false)
	if err == nil {
		t.Error("expected error for same source and destination")
	}
}

func TestCloneProjectEmptySource(t *testing.T) {
	s := newCloneStore()
	c := NewCloner(s)
	_, err := c.CloneProject("", "beta", false)
	if err == nil {
		t.Error("expected error for empty source")
	}
}

func TestCloneProjectEmptyDest(t *testing.T) {
	s := newCloneStore()
	c := NewCloner(s)
	_, err := c.CloneProject("alpha", "", false)
	if err == nil {
		t.Error("expected error for empty destination")
	}
}
