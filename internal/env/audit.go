package env

import (
	"fmt"
	"time"
)

// AuditEvent represents a recorded action on an environment variable.
type AuditEvent struct {
	Timestamp time.Time
	Project   string
	VarName   string
	Action    string // "set", "delete", "rotate", "copy", "move"
	Actor     string // os username or empty
}

// String returns a human-readable representation of the event.
func (e AuditEvent) String() string {
	return fmt.Sprintf("%s | %-8s | %s/%s | %s",
		e.Timestamp.Format(time.RFC3339),
		e.Action,
		e.Project,
		e.VarName,
		e.Actor,
	)
}

// AuditLog holds an ordered list of audit events.
type AuditLog struct {
	events []AuditEvent
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// Record appends a new event to the log.
func (l *AuditLog) Record(project, varName, action, actor string) {
	l.events = append(l.events, AuditEvent{
		Timestamp: time.Now().UTC(),
		Project:   project,
		VarName:   varName,
		Action:    action,
		Actor:     actor,
	})
}

// Events returns a copy of all recorded events.
func (l *AuditLog) Events() []AuditEvent {
	copy := make([]AuditEvent, len(l.events))
	for i, e := range l.events {
		copy[i] = e
	}
	return copy
}

// FilterByProject returns events matching the given project.
func (l *AuditLog) FilterByProject(project string) []AuditEvent {
	var result []AuditEvent
	for _, e := range l.events {
		if e.Project == project {
			result = append(result, e)
		}
	}
	return result
}

// FilterByAction returns events matching the given action.
func (l *AuditLog) FilterByAction(action string) []AuditEvent {
	var result []AuditEvent
	for _, e := range l.events {
		if e.Action == action {
			result = append(result, e)
		}
	}
	return result
}

// Len returns the number of recorded events.
func (l *AuditLog) Len() int {
	return len(l.events)
}
