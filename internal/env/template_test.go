package env

import (
	"testing"
)

func newTemplateStore(t *testing.T) *testStore {
	t.Helper()
	s := &testStore{data: map[string]map[string]string{}}
	s.data["myapp"] = map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"SECRET":  "s3cr3t",
	}
	return s
}

func TestRenderSimple(t *testing.T) {
	r, _ := NewTemplateRenderer(newTemplateStore(t), "myapp")
	out, err := r.Render("host={{DB_HOST}} port={{DB_PORT}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "host=localhost port=5432" {
		t.Errorf("got %q", out)
	}
}

func TestRenderMissingVar(t *testing.T) {
	r, _ := NewTemplateRenderer(newTemplateStore(t), "myapp")
	_, err := r.Render("key={{MISSING_KEY}}")
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestRenderNoPlaceholders(t *testing.T) {
	r, _ := NewTemplateRenderer(newTemplateStore(t), "myapp")
	out, err := r.Render("plain text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "plain text" {
		t.Errorf("got %q", out)
	}
}

func TestNewTemplateRendererEmptyProject(t *testing.T) {
	_, err := NewTemplateRenderer(newTemplateStore(t), "")
	if err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestListPlaceholders(t *testing.T) {
	names := ListPlaceholders("{{DB_HOST}}:{{DB_PORT}}/{{DB_HOST}}")
	if len(names) != 2 {
		t.Fatalf("expected 2 unique names, got %d", len(names))
	}
	if names[0] != "DB_HOST" || names[1] != "DB_PORT" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestListPlaceholdersEmpty(t *testing.T) {
	names := ListPlaceholders("no placeholders here")
	if len(names) != 0 {
		t.Errorf("expected empty, got %v", names)
	}
}

func TestRenderBytes(t *testing.T) {
	r, _ := NewTemplateRenderer(newTemplateStore(t), "myapp")
	out, err := r.RenderBytes([]byte("secret={{SECRET}}"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "secret=s3cr3t" {
		t.Errorf("got %q", string(out))
	}
}
