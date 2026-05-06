package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportFromDotenv(t *testing.T) {
	s := newTestStore(t)
	input := `
# comment
DB_HOST=localhost
DB_PORT="5432"
API_KEY=secret
`
	res, err := ImportFromDotenv(s, "myapp", strings.NewReader(input), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Imported) != 3 {
		t.Fatalf("expected 3 imported, got %d: %v", len(res.Imported), res.Imported)
	}
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	val, err := s.Get("myapp", "DB_PORT")
	if err != nil {
		t.Fatalf("get DB_PORT: %v", err)
	}
	if val != "5432" {
		t.Errorf("expected 5432, got %q", val)
	}
}

func TestImportSkipExisting(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("myapp", "DB_HOST", "original")
	input := "DB_HOST=new\nDB_PORT=5432\n"
	res, err := ImportFromDotenv(s, "myapp", strings.NewReader(input), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST skipped, got %v", res.Skipped)
	}
	if len(res.Imported) != 1 || res.Imported[0] != "DB_PORT" {
		t.Errorf("expected DB_PORT imported, got %v", res.Imported)
	}
	val, _ := s.Get("myapp", "DB_HOST")
	if val != "original" {
		t.Errorf("DB_HOST should remain original, got %q", val)
	}
}

func TestImportMalformedLine(t *testing.T) {
	s := newTestStore(t)
	input := "NOEQUALS\nVALID=ok\n"
	res, err := ImportFromDotenv(s, "myapp", strings.NewReader(input), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Errors) != 1 {
		t.Errorf("expected 1 error, got %v", res.Errors)
	}
	if len(res.Imported) != 1 {
		t.Errorf("expected 1 imported, got %v", res.Imported)
	}
}

func TestImportFromFile(t *testing.T) {
	s := newTestStore(t)
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("TOKEN=abc123\n"), 0600)
	res, err := ImportFromFile(s, "proj", path, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Imported) != 1 {
		t.Errorf("expected 1 imported, got %v", res.Imported)
	}
}

func TestImportFromFileMissing(t *testing.T) {
	s := newTestStore(t)
	_, err := ImportFromFile(s, "proj", "/nonexistent/.env", false)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestImportInvalidProject(t *testing.T) {
	s := newTestStore(t)
	_, err := ImportFromDotenv(s, "", strings.NewReader("K=V\n"), false)
	if err == nil {
		t.Error("expected error for empty project")
	}
}
