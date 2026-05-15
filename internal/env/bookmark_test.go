package env

import (
	"testing"
)

func newBookmarkStore() *BookmarkStore {
	return NewBookmarkStore(newMemKeychain())
}

func TestBookmarkSetAndResolve(t *testing.T) {
	s := newBookmarkStore()
	if err := s.Set("mydb", "production", "DB_PASSWORD"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	project, variable, err := s.Resolve("mydb")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if project != "production" || variable != "DB_PASSWORD" {
		t.Errorf("got %s/%s, want production/DB_PASSWORD", project, variable)
	}
}

func TestBookmarkResolveNotFound(t *testing.T) {
	s := newBookmarkStore()
	_, _, err := s.Resolve("missing")
	if err == nil {
		t.Fatal("expected error for missing bookmark")
	}
}

func TestBookmarkRemove(t *testing.T) {
	s := newBookmarkStore()
	_ = s.Set("tok", "proj", "API_KEY")
	if err := s.Remove("tok"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, _, err := s.Resolve("tok")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestBookmarkEmptyName(t *testing.T) {
	s := newBookmarkStore()
	if err := s.Set("", "proj", "VAR"); err == nil {
		t.Fatal("expected error for empty name")
	}
	if err := s.Remove(""); err == nil {
		t.Fatal("expected error for empty name on remove")
	}
	_, _, err := s.Resolve("")
	if err == nil {
		t.Fatal("expected error for empty name on resolve")
	}
}

func TestBookmarkList(t *testing.T) {
	s := newBookmarkStore()
	_ = s.Set("alpha", "proj1", "SECRET")
	_ = s.Set("beta", "proj2", "TOKEN")
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestBookmarkEntryString(t *testing.T) {
	e := BookmarkEntry{Name: "mydb", Project: "prod", Variable: "DB_PASS"}
	got := e.String()
	want := "mydb -> prod/DB_PASS"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBookmarkOverwrite(t *testing.T) {
	s := newBookmarkStore()
	_ = s.Set("key", "proj1", "OLD_VAR")
	_ = s.Set("key", "proj2", "NEW_VAR")
	project, variable, err := s.Resolve("key")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if project != "proj2" || variable != "NEW_VAR" {
		t.Errorf("got %s/%s, want proj2/NEW_VAR", project, variable)
	}
}
