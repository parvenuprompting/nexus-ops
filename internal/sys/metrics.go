package sys

import (
	"sync/atomic"
)

// Metrics holds thread-safe counters for the application.
// Metrics holds the app metrics.
// Note: We will likely send a DTO to frontend, but let's add tags just in case or create a DTO.
// Atomic fields don't marshal to JSON automatically.
// We should create a DTO struct for JSON output.

type Metrics struct {
	ActiveGoroutines atomic.Int64
	TasksCompleted   atomic.Int64
	Errors           atomic.Int64
}

type MetricsDTO struct {
	ActiveGoroutines int64 `json:"activeGoroutines"`
	TasksCompleted   int64 `json:"tasksCompleted"`
	Errors           int64 `json:"errors"`
}

// NewMetrics creates a new Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{}
}
