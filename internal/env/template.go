package env

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// templateVarRe matches {{VAR_NAME}} placeholders in templates.
var templateVarRe = regexp.MustCompile(`\{\{([A-Z_][A-Z0-9_]*)\}\}`)

// TemplateRenderer renders text templates by substituting environment
// variable placeholders with values loaded from the given Store.
type TemplateRenderer struct {
	store   Store
	project string
}

// NewTemplateRenderer returns a TemplateRenderer for the named project.
func NewTemplateRenderer(store Store, project string) (*TemplateRenderer, error) {
	if project == "" {
		return nil, fmt.Errorf("project name must not be empty")
	}
	return &TemplateRenderer{store: store, project: project}, nil
}

// Render replaces all {{VAR}} placeholders in src with the corresponding
// values from the keychain store. Missing variables are reported as errors.
func (r *TemplateRenderer) Render(src string) (string, error) {
	var missing []string
	var renderErr error

	result := templateVarRe.ReplaceAllStringFunc(src, func(match string) string {
		name := templateVarRe.FindStringSubmatch(match)[1]
		val, err := r.store.Get(r.project, name)
		if err != nil {
			missing = append(missing, name)
			return match
		}
		return val
	})

	if len(missing) > 0 {
		renderErr = fmt.Errorf("undefined variables: %s", strings.Join(missing, ", "))
	}
	return result, renderErr
}

// RenderBytes is a convenience wrapper around Render for byte slices.
func (r *TemplateRenderer) RenderBytes(src []byte) ([]byte, error) {
	out, err := r.Render(string(src))
	return []byte(out), err
}

// ListPlaceholders returns all unique variable names referenced in src.
func ListPlaceholders(src string) []string {
	seen := map[string]bool{}
	var names []string
	for _, m := range templateVarRe.FindAllStringSubmatch(src, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			names = append(names, m[1])
		}
	}
	return names
}

// RenderTo writes the rendered output to a bytes.Buffer.
func (r *TemplateRenderer) RenderTo(buf *bytes.Buffer, src string) error {
	out, err := r.Render(src)
	buf.WriteString(out)
	return err
}
