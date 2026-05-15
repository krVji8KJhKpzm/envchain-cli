package env

import (
	"errors"
	"fmt"
	"strings"
)

// NoteStore manages free-form notes attached to environment variables.
type NoteStore struct {
	kc keychain
}

func noteKey(project, varName string) string {
	return fmt.Sprintf("note::%s::%s", project, varName)
}

// NewNoteStore creates a NoteStore backed by the provided keychain.
func NewNoteStore(kc keychain) (*NoteStore, error) {
	if kc == nil {
		return nil, errors.New("note: keychain must not be nil")
	}
	return &NoteStore{kc: kc}, nil
}

// Set stores a note for the given project variable.
func (n *NoteStore) Set(project, varName, note string) error {
	if strings.TrimSpace(project) == "" {
		return errors.New("note: project must not be empty")
	}
	if strings.TrimSpace(varName) == "" {
		return errors.New("note: variable name must not be empty")
	}
	return n.kc.Set(noteKey(project, varName), note)
}

// Get retrieves the note for the given project variable.
// Returns an empty string and no error if no note is set.
func (n *NoteStore) Get(project, varName string) (string, error) {
	val, err := n.kc.Get(noteKey(project, varName))
	if err != nil {
		if isNotFound(err) {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

// Remove deletes the note for the given project variable.
func (n *NoteStore) Remove(project, varName string) error {
	err := n.kc.Delete(noteKey(project, varName))
	if isNotFound(err) {
		return nil
	}
	return err
}
