package todo

import "time"

// CreatedEvent is published after a todo is successfully persisted.
type CreatedEvent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}
