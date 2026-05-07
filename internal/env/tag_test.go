package env

import (
	"testing"
)

func newTagStore(t *testing.T) (*TagStore, Store) {
	t.Helper()
	s := newTestStore(t)
	ts, err := NewTagStore(s)
	if err != nil {
		t.Fatalf("NewTagStore: %v", err)
	}
	return ts, s
}

func TestAddAndListTags(t *testing.T) {
	ts, s := newTagStore(t)
	if err := s.Set("proj", "API_KEY", "secret"); err != nil {
		t.Fatal(err)
	}
	if err := ts.AddTag("proj", "API_KEY", "production"); err != nil {
		t.Fatal(err)
	}
	if err := ts.AddTag("proj", "API_KEY", "critical"); err != nil {
		t.Fatal(err)
	}
	tags, err := ts.ListTags("proj", "API_KEY")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "critical" || tags[1] != "production" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestRemoveTag(t *testing.T) {
	ts, _ := newTagStore(t)
	_ = ts.AddTag("proj", "DB_PASS", "sensitive")
	_ = ts.AddTag("proj", "DB_PASS", "backend")
	if err := ts.RemoveTag("proj", "DB_PASS", "sensitive"); err != nil {
		t.Fatal(err)
	}
	tags, _ := ts.ListTags("proj", "DB_PASS")
	if len(tags) != 1 || tags[0] != "backend" {
		t.Errorf("expected [backend], got %v", tags)
	}
}

func TestFindByTag(t *testing.T) {
	ts, _ := newTagStore(t)
	_ = ts.AddTag("proj", "API_KEY", "production")
	_ = ts.AddTag("proj", "DB_PASS", "production")
	_ = ts.AddTag("proj", "DEBUG", "dev")
	names, err := ts.FindByTag("proj", "production")
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 vars, got %d: %v", len(names), names)
	}
	if names[0] != "API_KEY" || names[1] != "DB_PASS" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestInvalidTag(t *testing.T) {
	ts, _ := newTagStore(t)
	if err := ts.AddTag("proj", "KEY", "bad tag!"); err == nil {
		t.Error("expected error for invalid tag, got nil")
	}
}

func TestNewTagStoreNilStore(t *testing.T) {
	_, err := NewTagStore(nil)
	if err == nil {
		t.Error("expected error for nil store")
	}
}
