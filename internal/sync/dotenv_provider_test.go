package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDotenvProviderPushAndPull(t *testing.T) {
	dir := t.TempDir()
	p := NewDotenvProvider(dir)

	vars := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "s3cr3t",
	}

	if err := p.Push("myproject", vars); err != nil {
		t.Fatalf("Push: %v", err)
	}

	got, err := p.Pull("myproject")
	if err != nil {
		t.Fatalf("Pull: %v", err)
	}
	for k, want := range vars {
		if got[k] != want {
			t.Errorf("var %s: want %q, got %q", k, want, got[k])
		}
	}
}

func TestDotenvProviderPullMissing(t *testing.T) {
	dir := t.TempDir()
	p := NewDotenvProvider(dir)

	vars, err := p.Pull("ghost")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(vars) != 0 {
		t.Errorf("expected empty map, got %v", vars)
	}
}

func TestDotenvProviderFilePermissions(t *testing.T) {
	dir := t.TempDir()
	p := NewDotenvProvider(dir)

	_ = p.Push("secure", map[string]string{"KEY": "val"})

	info, err := os.Stat(filepath.Join(dir, "secure.env"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected 0600 permissions, got %o", perm)
	}
}

func TestDotenvProviderName(t *testing.T) {
	p := NewDotenvProvider("/tmp")
	if p.Name() != DotenvProviderName {
		t.Errorf("expected %q, got %q", DotenvProviderName, p.Name())
	}
}
