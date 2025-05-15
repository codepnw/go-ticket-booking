package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, e *dto.EventRequest) (*domain.Event, error)
	ListEvents(ctx context.Context) ([]*domain.Event, error)
	GetEventByID(ctx context.Context, id int) (*domain.Event, error)
	UpdateEvent(ctx context.Context, id int, e *dto.EventRequest) error
	DeleteEvent(ctx context.Context, id int) error
}

type eventUsecase struct {
	repo repository.EventRepository
}

func NewEventUsecase(repo repository.EventRepository) EventUsecase {
	return &eventUsecase{repo: repo}
}

func (u *eventUsecase) CreateEvent(ctx context.Context, req *dto.EventRequest) (*domain.Event, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if req.StartTime == nil || req.EndTime == nil {
		return nil, errors.New("start and end time is required")
	}

	if req.EndTime.Before(*req.StartTime) {
		return nil, errors.New("end time cannot before start ime")
	}

	event := &domain.Event{
		Name:        req.Name,
		Description: req.Description,
		StartTime:   *req.StartTime,
		EndTime:     *req.EndTime,
		// TODO: update later
		LocationID: 1,
	}

	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	if err := u.repo.Create(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (u *eventUsecase) UpdateEvent(ctx context.Context, id int, req *dto.EventRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	exist, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Name == "" {
		req.Name = exist.Name
	}

	if req.Description == "" {
		req.Description = exist.Description
	}

	if req.StartTime == nil {
		*req.StartTime = exist.StartTime
	}

	if req.EndTime == nil {
		*req.StartTime = exist.EndTime
	}

	if req.LocationID == 0 {
		req.LocationID = exist.LocationID
	}

	if req.EndTime.Before(*req.StartTime) {
		return errors.New("end time cannot before start ime")
	}

	exist.UpdatedAt = time.Now()

	event := domain.Event{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		StartTime:   *req.StartTime,
		EndTime:     *req.EndTime,
		LocationID:  req.LocationID,
		UpdatedAt:   exist.UpdatedAt,
	}

	return u.repo.Update(ctx, &event)
}

func (u *eventUsecase) ListEvents(ctx context.Context) ([]*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.List(ctx)
}

func (u *eventUsecase) GetEventByID(ctx context.Context, id int) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *eventUsecase) DeleteEvent(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
