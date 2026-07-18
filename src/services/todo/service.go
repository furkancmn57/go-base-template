package todo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/furkancmn57/go-base-template/src/common/apperr"
	"github.com/furkancmn57/go-base-template/src/constants"
	"github.com/furkancmn57/go-base-template/src/data/entities"
	todoevents "github.com/furkancmn57/go-base-template/src/events/todo"
	"github.com/furkancmn57/go-base-template/src/interfaces"
	"github.com/furkancmn57/go-base-template/src/models/requests"
	"github.com/furkancmn57/go-base-template/src/models/responses"
	"github.com/furkancmn57/go-base-template/src/services/todo/validations"
)

// Service implements every todo use-case directly against *gorm.DB —
// there is no repository layer in this architecture.
type Service struct {
	db        *gorm.DB
	publisher interfaces.Publisher
}

// NewService wires a Service with its GORM handle and event publisher.
func NewService(db *gorm.DB, publisher interfaces.Publisher) *Service {
	return &Service{db: db, publisher: publisher}
}

// Create validates and persists a new todo, then best-effort publishes
// todo.created.
func (s *Service) Create(ctx context.Context, req requests.CreateTodoRequest) (*responses.TodoResponse, *apperr.Error) {
	if err := validations.CreateTodoRequest(req); err != nil {
		return nil, err
	}

	entity := entities.Todo{Title: req.Title, Description: req.Description}
	if err := s.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, apperr.Internal(err)
	}

	s.publish(ctx, constants.TodoCreated, todoevents.CreatedEvent{
		ID:        entity.ID.String(),
		Title:     entity.Title,
		CreatedAt: entity.CreatedAt,
	})

	resp := toResponse(entity)
	return &resp, nil
}

// Todos lists every todo, most recent first.
func (s *Service) Todos(ctx context.Context) ([]responses.TodoResponse, *apperr.Error) {
	var rows []entities.Todo
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, apperr.Internal(err)
	}

	out := make([]responses.TodoResponse, 0, len(rows))
	for _, entity := range rows {
		out = append(out, toResponse(entity))
	}
	return out, nil
}

// TodoById fetches a single todo by its ID.
func (s *Service) TodoById(ctx context.Context, id string) (*responses.TodoResponse, *apperr.Error) {
	entity, appErr := s.findByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}
	resp := toResponse(*entity)
	return &resp, nil
}

// Update validates and applies changes to an existing todo, then best-effort
// publishes todo.updated.
func (s *Service) Update(ctx context.Context, id string, req requests.UpdateTodoRequest) (*responses.TodoResponse, *apperr.Error) {
	if err := validations.UpdateTodoRequest(req); err != nil {
		return nil, err
	}

	entity, appErr := s.findByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	entity.Title = req.Title
	entity.Description = req.Description
	entity.Completed = req.Completed

	if err := s.db.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, apperr.Internal(err)
	}

	s.publish(ctx, constants.TodoUpdated, todoevents.UpdatedEvent{
		ID:        entity.ID.String(),
		Title:     entity.Title,
		UpdatedAt: entity.UpdatedAt,
	})

	resp := toResponse(*entity)
	return &resp, nil
}

// Complete marks a todo as completed, then best-effort publishes
// todo.completed.
func (s *Service) Complete(ctx context.Context, id string) (*responses.TodoResponse, *apperr.Error) {
	entity, appErr := s.findByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	if entity.Completed {
		return nil, apperr.Conflict(constants.TodoAlreadyCompleted, "todo is already completed")
	}

	entity.Completed = true
	if err := s.db.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, apperr.Internal(err)
	}

	s.publish(ctx, constants.TodoCompleted, todoevents.CompletedEvent{
		ID:          entity.ID.String(),
		CompletedAt: entity.UpdatedAt,
	})

	resp := toResponse(*entity)
	return &resp, nil
}

// Delete soft-deletes a todo, then best-effort publishes todo.deleted.
func (s *Service) Delete(ctx context.Context, id string) *apperr.Error {
	entity, appErr := s.findByID(ctx, id)
	if appErr != nil {
		return appErr
	}

	if err := s.db.WithContext(ctx).Delete(entity).Error; err != nil {
		return apperr.Internal(err)
	}

	s.publish(ctx, constants.TodoDeleted, todoevents.DeletedEvent{
		ID:        entity.ID.String(),
		DeletedAt: time.Now(),
	})

	return nil
}

func (s *Service) findByID(ctx context.Context, id string) (*entities.Todo, *apperr.Error) {
	todoID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperr.BadRequest(constants.TodoInvalidID, "invalid todo id")
	}

	var entity entities.Todo
	if err := s.db.WithContext(ctx).First(&entity, "id = ?", todoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound(constants.TodoNotFound, "todo not found")
		}
		return nil, apperr.Internal(err)
	}
	return &entity, nil
}

// publish fires a domain event on a best-effort basis: failures are logged
// and never roll back the DB transaction that already committed.
func (s *Service) publish(ctx context.Context, topic string, payload any) {
	if s.publisher == nil {
		return
	}
	if err := s.publisher.Publish(ctx, topic, payload); err != nil {
		log.Printf("todo: best-effort publish failed for topic %q: %v", topic, err)
	}
}

func toResponse(entity entities.Todo) responses.TodoResponse {
	return responses.TodoResponse{
		ID:          entity.ID.String(),
		Title:       entity.Title,
		Description: entity.Description,
		Completed:   entity.Completed,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
