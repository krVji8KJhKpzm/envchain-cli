package env

import (
	"testing"
)

func newIntegrationBookmarkStore() *BookmarkStore {
	return NewBookmarkStore(newMemKeychain())
}

func TestBookmarkIntegrationSetListResolve(t *testing.T) {
	s := newIntegrationBookmarkStore()

	bookmarks := []struct {
		name, project, variable string
	}{
		{"db", "production", "DB_PASSWORD"},
		{"api", "staging", "API_TOKEN"},
		{"cache", "production", "REDIS_URL"},
	}

	for _, bm := range bookmarks {
		if err := s.Set(bm.name, bm.project, bm.variable); err != nil {
			t.Fatalf("Set(%s): %v", bm.name, err)
		}
	}

	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != len(bookmarks) {
		t.Fatalf("expected %d entries, got %d", len(bookmarks), len(entries))
	}

	for _, bm := range bookmarks {
		p, v, err := s.Resolve(bm.name)
		if err != nil {
			t.Errorf("Resolve(%s): %v", bm.name, err)
			continue
		}
		if p != bm.project || v != bm.variable {
			t.Errorf("Resolve(%s) = %s/%s, want %s/%s", bm.name, p, v, bm.project, bm.variable)
		}
	}
}

func TestBookmarkIntegrationRemoveReducesList(t *testing.T) {
	s := newIntegrationBookmarkStore()
	_ = s.Set("first", "proj", "VAR1")
	_ = s.Set("second", "proj", "VAR2")

	if err := s.Remove("first"); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after remove, got %d", len(entries))
	}
	if entries[0].Name != "second" {
		t.Errorf("expected remaining entry to be 'second', got %q", entries[0].Name)
	}
}
