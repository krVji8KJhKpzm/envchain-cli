package env

import (
	"encoding/json"
	"fmt"
	"strings"
)

const defaultQuotaLimit = 100

// QuotaStore manages per-project variable count limits.
type QuotaStore struct {
	kc keychain
	project string
}

// QuotaEntry holds the limit and current usage for a project.
type QuotaEntry struct {
	Limit int `json:"limit"`
}

// NewQuotaStore creates a QuotaStore for the given project.
func NewQuotaStore(kc keychain, project string) (*QuotaStore, error) {
	if strings.TrimSpace(project) == "" {
		return nil, fmt.Errorf("project name must not be empty")
	}
	return &QuotaStore{kc: kc, project: project}, nil
}

func quotaKey(project string) string {
	return "quota:" + project
}

// SetLimit sets the maximum number of variables allowed for the project.
func (q *QuotaStore) SetLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("quota limit must be greater than zero")
	}
	entry := QuotaEntry{Limit: limit}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal quota: %w", err)
	}
	return q.kc.Set(q.project, quotaKey(q.project), string(data))
}

// GetLimit returns the configured limit, or the default if none is set.
func (q *QuotaStore) GetLimit() (int, error) {
	raw, err := q.kc.Get(q.project, quotaKey(q.project))
	if err != nil {
		return defaultQuotaLimit, nil
	}
	var entry QuotaEntry
	if err := json.Unmarshal([]byte(raw), &entry); err != nil {
		return 0, fmt.Errorf("unmarshal quota: %w", err)
	}
	return entry.Limit, nil
}

// CheckQuota returns an error if adding delta more variables would exceed the limit.
func (q *QuotaStore) CheckQuota(current, delta int) error {
	limit, err := q.GetLimit()
	if err != nil {
		return err
	}
	if current+delta > limit {
		return fmt.Errorf("quota exceeded: project %q allows %d variables, currently at %d",
			q.project, limit, current)
	}
	return nil
}

// RemoveLimit deletes the quota entry, reverting to the default.
func (q *QuotaStore) RemoveLimit() error {
	return q.kc.Delete(q.project, quotaKey(q.project))
}
