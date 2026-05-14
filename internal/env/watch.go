// Package env provides environment variable management backed by a keychain.
package env

import (
	"fmt"
	"sync"
	"time"
)

// WatchEvent describes a change detected for a watched variable.
type WatchEvent struct {
	Project  string
	VarName  string
	OldValue string
	NewValue string
	At       time.Time
}

func (e WatchEvent) String() string {
	return fmt.Sprintf("[%s] %s/%s changed at %s",
		e.Project, e.VarName, e.VarName, e.At.Format(time.RFC3339))
}

// Watcher polls a Store for changes to registered variables and emits events.
type Watcher struct {
	store    Store
	interval time.Duration
	watched  map[string]map[string]string // project -> varName -> lastValue
	mu       sync.Mutex
	stopCh   chan struct{}
}

// NewWatcher creates a Watcher that polls at the given interval.
func NewWatcher(store Store, interval time.Duration) *Watcher {
	return &Watcher{
		store:    store,
		interval: interval,
		watched:  make(map[string]map[string]string),
		stopCh:   make(chan struct{}),
	}
}

// Add registers a variable for watching within a project.
func (w *Watcher) Add(project, varName string) error {
	if project == "" {
		return fmt.Errorf("project must not be empty")
	}
	if varName == "" {
		return fmt.Errorf("varName must not be empty")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.watched[project] == nil {
		w.watched[project] = make(map[string]string)
	}
	val, _ := w.store.Get(project, varName)
	w.watched[project][varName] = val
	return nil
}

// Remove unregisters a variable from watching.
func (w *Watcher) Remove(project, varName string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if m, ok := w.watched[project]; ok {
		delete(m, varName)
	}
}

// Start begins polling and sends events to the returned channel.
// Call Stop to terminate the watcher.
func (w *Watcher) Start() <-chan WatchEvent {
	ch := make(chan WatchEvent, 16)
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		defer close(ch)
		for {
			select {
			case <-w.stopCh:
				return
			case <-ticker.C:
				w.poll(ch)
			}
		}
	}()
	return ch
}

// Stop terminates the polling loop.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

func (w *Watcher) poll(ch chan<- WatchEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for project, vars := range w.watched {
		for varName, last := range vars {
			current, _ := w.store.Get(project, varName)
			if current != last {
				ch <- WatchEvent{
					Project:  project,
					VarName:  varName,
					OldValue: last,
					NewValue: current,
					At:       time.Now(),
				}
				w.watched[project][varName] = current
			}
		}
	}
}
