package sys

import (
	"context"
	"sync"
)

type AppState struct {
	Ctx     context.Context
	Cancel  context.CancelFunc
	Metrics *Metrics

	// Phase 1: Stability - Replace sync.Map with Type-safe Map + Mutex
	activeTasksMu sync.RWMutex
	activeTasks   map[string]*Task

	// Phase 2: UI Responsiveness - Update Channel
	UpdateChan chan struct{}
}

// NewAppState creates a new state with initialized maps and channels
func NewAppState(ctx context.Context, cancel context.CancelFunc, metrics *Metrics) *AppState {
	return &AppState{
		Ctx:         ctx,
		Cancel:      cancel,
		Metrics:     metrics,
		activeTasks: make(map[string]*Task),
		UpdateChan:  make(chan struct{}, 1), // Buffered 1 for debouncing
	}
}

// NotifyUpdate sends a non-blocking signal to the UI
func (s *AppState) NotifyUpdate() {
	select {
	case s.UpdateChan <- struct{}{}:
	default:
		// Channel full, UI has pending update, skip
	}
}

// AddActiveTask adds a task and notifies listeners
func (s *AppState) AddActiveTask(t *Task) {
	s.activeTasksMu.Lock()
	s.activeTasks[t.ID] = t
	s.activeTasksMu.Unlock()
	s.NotifyUpdate()
}

// RemoveActiveTask removes a task and notifies listeners
func (s *AppState) RemoveActiveTask(id string) {
	s.activeTasksMu.Lock()
	delete(s.activeTasks, id)
	s.activeTasksMu.Unlock()
	s.NotifyUpdate()
}

// SnapshotActiveTasks returns a safe copy of all active tasks
func (s *AppState) SnapshotActiveTasks() []*Task {
	s.activeTasksMu.RLock()
	defer s.activeTasksMu.RUnlock()

	tasks := make([]*Task, 0, len(s.activeTasks))
	for _, t := range s.activeTasks {
		tasks = append(tasks, t)
	}
	return tasks
}
