package requests

// UpdateTodoRequest is the payload for PUT /todos/:id.
type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
