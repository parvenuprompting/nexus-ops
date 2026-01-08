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
	ID        string    `json:"id"`
	Type      TaskType  `json:"type"`
	StartedAt time.Time `json:"startedAt"`
	Status    string    `json:"status"`
}
