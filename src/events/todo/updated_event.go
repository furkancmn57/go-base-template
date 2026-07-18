package todo

import "time"

// UpdatedEvent is published after a todo is successfully updated.
type UpdatedEvent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
}
