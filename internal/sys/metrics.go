package sys

import (
	"sync/atomic"
)

// Metrics holds thread-safe counters for the application.
type Metrics struct {
	ActiveGoroutines atomic.Int64
	TasksCompleted   atomic.Int64
	Errors           atomic.Int64
}

// NewMetrics creates a new Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{}
}
