package todo

import (
	"context"
	"encoding/json"
	"log"

	"github.com/furkancmn57/go-base-template/src/common"
	"github.com/furkancmn57/go-base-template/src/constants"
	todoevents "github.com/furkancmn57/go-base-template/src/events/todo"
	"github.com/furkancmn57/go-base-template/src/interfaces"
)

// LogWhenTodoCompleted reacts to todo.completed (wiring demo).
type LogWhenTodoCompleted struct{}

// NewLogWhenTodoCompleted constructs the consumer.
func NewLogWhenTodoCompleted() *LogWhenTodoCompleted {
	return &LogWhenTodoCompleted{}
}

// Register subscribes this consumer to todo.completed.
func (c *LogWhenTodoCompleted) Register(subscriber interfaces.Subscriber) error {
	handler := interfaces.Handler(common.BaseConsumer("log-when-todo-completed", c.consume))
	return subscriber.Subscribe(constants.TodoCompleted, handler)
}

func (c *LogWhenTodoCompleted) consume(_ context.Context, payload []byte) error {
	var event todoevents.CompletedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	log.Printf("todo: completed id=%s at=%s", event.ID, event.CompletedAt)
	return nil
}
