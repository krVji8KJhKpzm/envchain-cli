package env

import (
	"testing"
	"time"
)

// watchStore is an in-memory Store for watcher tests.
type watchStore struct {
	data map[string]map[string]string
}

func newWatchStore() *watchStore {
	return &watchStore{data: make(map[string]map[string]string)}
}

func (s *watchStore) Set(project, key, value string) error {
	if s.data[project] == nil {
		s.data[project] = make(map[string]string)
	}
	s.data[project][key] = value
	return nil
}

func (s *watchStore) Get(project, key string) (string, error) {
	if m, ok := s.data[project]; ok {
		if v, ok2 := m[key]; ok2 {
			return v, nil
		}
	}
	return "", ErrNotFound
}

func (s *watchStore) Delete(project, key string) error {
	if m, ok := s.data[project]; ok {
		delete(m, key)
	}
	return nil
}

func TestWatcherDetectsChange(t *testing.T) {
	st := newWatchStore()
	_ = st.Set("proj", "TOKEN", "old")

	w := NewWatcher(st, 20*time.Millisecond)
	if err := w.Add("proj", "TOKEN"); err != nil {
		t.Fatalf("Add: %v", err)
	}

	ch := w.Start()
	defer w.Stop()

	_ = st.Set("proj", "TOKEN", "new")

	select {
	case ev := <-ch:
		if ev.OldValue != "old" || ev.NewValue != "new" {
			t.Errorf("expected old=old new=new, got %q %q", ev.OldValue, ev.NewValue)
		}
		if ev.Project != "proj" || ev.VarName != "TOKEN" {
			t.Errorf("unexpected event fields: %+v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for watch event")
	}
}

func TestWatcherNoEventWhenUnchanged(t *testing.T) {
	st := newWatchStore()
	_ = st.Set("proj", "KEY", "same")

	w := NewWatcher(st, 20*time.Millisecond)
	_ = w.Add("proj", "KEY")
	ch := w.Start()

	select {
	case ev := <-ch:
		t.Errorf("unexpected event: %+v", ev)
	case <-time.After(100 * time.Millisecond):
		// expected — no change
	}
	w.Stop()
}

func TestWatcherAddEmptyProject(t *testing.T) {
	st := newWatchStore()
	w := NewWatcher(st, time.Second)
	if err := w.Add("", "KEY"); err == nil {
		t.Error("expected error for empty project")
	}
}

func TestWatcherAddEmptyVarName(t *testing.T) {
	st := newWatchStore()
	w := NewWatcher(st, time.Second)
	if err := w.Add("proj", ""); err == nil {
		t.Error("expected error for empty varName")
	}
}

func TestWatcherRemove(t *testing.T) {
	st := newWatchStore()
	_ = st.Set("proj", "A", "v1")

	w := NewWatcher(st, 20*time.Millisecond)
	_ = w.Add("proj", "A")
	w.Remove("proj", "A")

	ch := w.Start()
	_ = st.Set("proj", "A", "v2")

	select {
	case ev := <-ch:
		t.Errorf("unexpected event after Remove: %+v", ev)
	case <-time.After(100 * time.Millisecond):
		// expected
	}
	w.Stop()
}

func TestWatchEventString(t *testing.T) {
	ev := WatchEvent{Project: "p", VarName: "V", OldValue: "a", NewValue: "b", At: time.Time{}}
	s := ev.String()
	if s == "" {
		t.Error("WatchEvent.String() returned empty string")
	}
}
