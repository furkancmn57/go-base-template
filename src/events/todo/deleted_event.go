package todo

import "time"

// DeletedEvent is published after a todo is soft-deleted.
type DeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deletedAt"`
}
