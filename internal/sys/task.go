package sys

import "time"

// TaskType represents the category of a task (CPU, IO, etc.)
type TaskType string

const (
	TaskTypeCPU         TaskType = "CPU-Bound"
	TaskTypeIO          TaskType = "I/O-Bound"
	TaskTypeImageResize TaskType = "Image-Resize"
)

type Task struct {
	ID        string
	Type      TaskType
	StartedAt time.Time
	Status    string
}
