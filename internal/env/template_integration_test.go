package env_test

import (
	"testing"

	"github.com/envchain/envchain-cli/internal/env"
)

// integrationStore is a simple in-memory Store for integration tests.
type integrationStore struct {
	data map[string]map[string]string
}

func (s *integrationStore) Get(project, name string) (string, error) {
	if p, ok := s.data[project]; ok {
		if v, ok := p[name]; ok {
			return v, nil
		}
	}
	return "", env.ErrNotFound
}

func (s *integrationStore) Set(project, name, value string) error {
	if s.data[project] == nil {
		s.data[project] = map[string]string{}
	}
	s.data[project][name] = value
	return nil
}

func (s *integrationStore) Delete(project, name string) error {
	delete(s.data[project], name)
	return nil
}

func TestTemplateRenderIntegration(t *testing.T) {
	store := &integrationStore{
		data: map[string]map[string]string{
			"prod": {
				"API_URL":  "https://api.example.com",
				"API_KEY":  "abc123",
				"TIMEOUT": "30",
			},
		},
	}

	r, err := env.NewTemplateRenderer(store, "prod")
	if err != nil {
		t.Fatalf("NewTemplateRenderer: %v", err)
	}

	tmpl := "curl -H 'Authorization: {{API_KEY}}' {{API_URL}}/health --max-time {{TIMEOUT}}"
	got, err := r.Render(tmpl)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	want := "curl -H 'Authorization: abc123' https://api.example.com/health --max-time 30"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestTemplatePartialMissingIntegration(t *testing.T) {
	store := &integrationStore{
		data: map[string]map[string]string{
			"staging": {"HOST": "staging.example.com"},
		},
	}
	r, _ := env.NewTemplateRenderer(store, "staging")
	_, err := r.Render("{{HOST}}:{{PORT}}")
	if err == nil {
		t.Fatal("expected error for missing PORT")
	}
}
