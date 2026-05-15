package env

import "fmt"

const readonlyPrefix = "readonly:"

// ReadonlyStore manages read-only flags for environment variables.
// A read-only variable can be read but not overwritten or deleted without
// explicitly removing the read-only flag first.
type ReadonlyStore struct {
	kc keychain
	project string
}

// NewReadonlyStore creates a ReadonlyStore for the given project.
func NewReadonlyStore(kc keychain, project string) (*ReadonlyStore, error) {
	if project == "" {
		return nil, fmt.Errorf("project name must not be empty")
	}
	return &ReadonlyStore{kc: kc, project: project}, nil
}

func readonlyKey(project, varName string) string {
	return readonlyPrefix + project + ":" + varName
}

// SetReadonly marks varName as read-only.
func (r *ReadonlyStore) SetReadonly(varName string) error {
	if varName == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	return r.kc.Set(readonlyKey(r.project, varName), "true")
}

// IsReadonly reports whether varName is currently marked read-only.
func (r *ReadonlyStore) IsReadonly(varName string) (bool, error) {
	val, err := r.kc.Get(readonlyKey(r.project, varName))
	if err != nil {
		if err.Error() == "secret not found" {
			return false, nil
		}
		return false, err
	}
	return val == "true", nil
}

// Unset removes the read-only flag from varName.
func (r *ReadonlyStore) Unset(varName string) error {
	if varName == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	err := r.kc.Delete(readonlyKey(r.project, varName))
	if err != nil && err.Error() == "secret not found" {
		return fmt.Errorf("readonly flag not set for %q", varName)
	}
	return err
}

// ListReadonly returns all variable names that are marked read-only in the project.
func (r *ReadonlyStore) ListReadonly() ([]string, error) {
	prefix := readonlyPrefix + r.project + ":"
	keys, err := r.kc.List(prefix)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(keys))
	for _, k := range keys {
		names = append(names, k[len(prefix):])
	}
	return names, nil
}
