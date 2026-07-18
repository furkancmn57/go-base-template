package requests

// CreateTodoRequest is the payload for POST /todos.
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
