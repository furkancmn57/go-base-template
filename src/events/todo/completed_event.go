package todo

import "time"

// CompletedEvent is published after a todo is marked as completed.
type CompletedEvent struct {
	ID          string    `json:"id"`
	CompletedAt time.Time `json:"completedAt"`
}
