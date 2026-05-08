package env

import (
	"fmt"
	"regexp"
	"strings"
)

var validNamespaceRe = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

// Namespace represents a logical grouping of projects under a common prefix.
type Namespace struct {
	Name string
}

// NamespaceStore manages namespace-to-project mappings via the underlying Store.
type NamespaceStore struct {
	store Store
}

// NewNamespaceStore creates a NamespaceStore backed by the given Store.
func NewNamespaceStore(s Store) (*NamespaceStore, error) {
	if s == nil {
		return nil, fmt.Errorf("store must not be nil")
	}
	return &NamespaceStore{store: s}, nil
}

// ProjectKey returns the fully-qualified project key for a namespace and project.
func (ns *NamespaceStore) ProjectKey(namespace, project string) (string, error) {
	if err := validateNamespace(namespace); err != nil {
		return "", err
	}
	if strings.TrimSpace(project) == "" {
		return "", fmt.Errorf("project must not be empty")
	}
	return namespace + "/" + project, nil
}

// ListProjects returns all projects that belong to the given namespace.
func (ns *NamespaceStore) ListProjects(namespace string, allProjects []string) ([]string, error) {
	if err := validateNamespace(namespace); err != nil {
		return nil, err
	}
	prefix := namespace + "/"
	var result []string
	for _, p := range allProjects {
		if strings.HasPrefix(p, prefix) {
			result = append(result, strings.TrimPrefix(p, prefix))
		}
	}
	return result, nil
}

// ParseProjectKey splits a fully-qualified project key into namespace and project.
// Returns an error if the key does not contain a namespace prefix.
func ParseProjectKey(key string) (namespace, project string, err error) {
	parts := strings.SplitN(key, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid project key %q: expected <namespace>/<project>", key)
	}
	return parts[0], parts[1], nil
}

func validateNamespace(namespace string) error {
	if strings.TrimSpace(namespace) == "" {
		return fmt.Errorf("namespace must not be empty")
	}
	if !validNamespaceRe.MatchString(namespace) {
		return fmt.Errorf("namespace %q contains invalid characters", namespace)
	}
	return nil
}
