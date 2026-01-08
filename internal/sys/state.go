package sys

import (
	"context"
	"sync"
)

// AppState holds the global application state and dependencies.
type AppState struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Metrics *Metrics
	ActiveTasks sync.Map // map[string]*Task
}
