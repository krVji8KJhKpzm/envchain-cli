package env

import (
	"encoding/json"
	"fmt"
	"time"
)

// HistoryEntry records a previous value of an environment variable.
type HistoryEntry struct {
	Value     string    `json:"value"`
	ChangedAt time.Time `json:"changed_at"`
	Actor     string    `json:"actor,omitempty"`
}

// HistoryStore tracks value history for environment variables.
type HistoryStore struct {
	store   Store
	maxSize int
}

// NewHistoryStore creates a HistoryStore backed by the given Store.
// maxSize controls how many historical entries are retained per variable.
func NewHistoryStore(s Store, maxSize int) *HistoryStore {
	if maxSize <= 0 {
		maxSize = 10
	}
	return &HistoryStore{store: s, maxSize: maxSize}
}

func historyKey(project, name string) string {
	return fmt.Sprintf("__history__%s__%s", project, name)
}

// Record appends a new history entry for the given project/variable.
func (h *HistoryStore) Record(project, name, value, actor string) error {
	entries, err := h.List(project, name)
	if err != nil {
		entries = []HistoryEntry{}
	}

	entries = append([]HistoryEntry{{
		Value:     value,
		ChangedAt: time.Now().UTC(),
		Actor:     actor,
	}}, entries...)

	if len(entries) > h.maxSize {
		entries = entries[:h.maxSize]
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	return h.store.Set(project, historyKey(project, name), string(data))
}

// List returns the recorded history for a variable, newest first.
func (h *HistoryStore) List(project, name string) ([]HistoryEntry, error) {
	raw, err := h.store.Get(project, historyKey(project, name))
	if err != nil {
		return nil, fmt.Errorf("history: not found: %w", err)
	}
	var entries []HistoryEntry
	if err := json.Unmarshal([]byte(raw), &entries); err != nil {
		return nil, fmt.Errorf("history: unmarshal: %w", err)
	}
	return entries, nil
}

// Clear removes all recorded history for a variable.
func (h *HistoryStore) Clear(project, name string) error {
	return h.store.Delete(project, historyKey(project, name))
}
