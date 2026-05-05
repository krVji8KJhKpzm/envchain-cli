package sync

import (
	"fmt"
)

// Provider defines the interface for team sync backends.
type Provider interface {
	// Push uploads the given key-value pairs for the specified project.
	Push(project string, vars map[string]string) error
	// Pull retrieves key-value pairs for the specified project.
	Pull(project string) (map[string]string, error)
	// Name returns the provider identifier.
	Name() string
}

// ErrProviderNotFound is returned when a requested provider is not registered.
type ErrProviderNotFound struct {
	Name string
}

func (e *ErrProviderNotFound) Error() string {
	return fmt.Sprintf("sync provider %q not found", e.Name)
}

var registry = map[string]Provider{}

// Register adds a provider to the global registry.
func Register(p Provider) {
	registry[p.Name()] = p
}

// Get retrieves a registered provider by name.
func Get(name string) (Provider, error) {
	p, ok := registry[name]
	if !ok {
		return nil, &ErrProviderNotFound{Name: name}
	}
	return p, nil
}

// Available returns the names of all registered providers.
func Available() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	return names
}
