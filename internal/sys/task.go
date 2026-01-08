package sys

import (
	"time"
)

type TaskType string

const (
	TaskTypeCPU TaskType = "CPU-Bound"
	TaskTypeIO  TaskType = "I/O-Bound"
)

// Task represents a unit of work.
type Task struct {
	ID        string
	Type      TaskType
	StartedAt time.Time
	Status    string
}
