package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain-cli/internal/sync"
)

func TestNewManifest(t *testing.T) {
	m := sync.NewManifest("myproject")
	if m.Project != "myproject" {
		t.Fatalf("expected project 'myproject', got %q", m.Project)
	}
	if m.Version != 1 {
		t.Fatalf("expected version 1, got %d", m.Version)
	}
	if m.Vars == nil {
		t.Fatal("expected non-nil Vars map")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	m := sync.NewManifest("testproject")
	m.Vars["DATABASE_URL"] = sync.VarMeta{Description: "Postgres connection string", Required: true}
	m.Vars["API_KEY"] = sync.VarMeta{Required: false}

	if err := m.Save(dir); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := sync.LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest() error: %v", err)
	}
	if loaded.Project != "testproject" {
		t.Errorf("expected project 'testproject', got %q", loaded.Project)
	}
	if len(loaded.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(loaded.Vars))
	}
	if !loaded.Vars["DATABASE_URL"].Required {
		t.Error("expected DATABASE_URL to be required")
	}
}

func TestLoadManifestNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := sync.LoadManifest(dir)
	if err != sync.ErrManifestNotFound {
		t.Fatalf("expected ErrManifestNotFound, got %v", err)
	}
}

func TestSaveCreatesFile(t *testing.T) {
	dir := t.TempDir()
	m := sync.NewManifest("proj")
	if err := m.Save(dir); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	path := filepath.Join(dir, sync.ManifestFileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected manifest file to exist at %s", path)
	}
}

func TestSaveUpdatesTimestamp(t *testing.T) {
	dir := t.TempDir()
	m := sync.NewManifest("proj")
	if err := m.Save(dir); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	if m.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set after Save")
	}
}
