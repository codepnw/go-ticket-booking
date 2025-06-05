package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/errs"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, e *dto.EventRequest) (*domain.Event, error)
	ListEvents(ctx context.Context) ([]*domain.Event, error)
	GetEventByID(ctx context.Context, id int64) (*domain.Event, error)
	UpdateEvent(ctx context.Context, id int64, e *dto.EventUpdateRequest) error
	DeleteEvent(ctx context.Context, id int64) error

	// Location
	CreateLocation(ctx context.Context, req *dto.LocationRequest) error
	ListLocations(ctx context.Context, limit, offset int) ([]*domain.Location, error)
	GetLocationByID(ctx context.Context, id int64) (*domain.Location, error)
	UpdateLocation(ctx context.Context, id int64, req *dto.LocationUpdateRequest) error
	DeleteLocation(ctx context.Context, id int64) error
}

type eventUsecase struct {
	repo repository.EventRepository
}

func NewEventUsecase(repo repository.EventRepository) EventUsecase {
	return &eventUsecase{
		repo: repo,
	}
}

func (u *eventUsecase) CreateEvent(ctx context.Context, req *dto.EventRequest) (*domain.Event, error) {
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time cannot before start time")
	}

	event := &domain.Event{
		Name:        req.Name,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		LocationID:  req.LocationID,
	}

	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	if err := u.repo.CreateEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func (u *eventUsecase) UpdateEvent(ctx context.Context, id int64, req *dto.EventUpdateRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	event, err := u.repo.GetEventByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Name == nil && req.Description == nil && req.StartTime == nil &&
		req.EndTime == nil && req.LocationID == nil {
		return errs.ErrNoFieldsToUpdate
	}

	if req.Name != nil {
		event.Name = *req.Name
	}

	if req.Description != nil {
		event.Description = *req.Description
	}

	if req.StartTime != nil {
		event.StartTime = *req.StartTime
	}

	if req.EndTime != nil {
		event.EndTime = *req.EndTime
	}

	if req.LocationID != nil {
		event.LocationID = *req.LocationID
	}

	if event.EndTime.Before(event.StartTime) {
		return errors.New("end time cannot be before start time")
	}

	event.UpdatedAt = time.Now()

	return u.repo.UpdateEvent(ctx, event)
}

func (u *eventUsecase) ListEvents(ctx context.Context) ([]*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.ListEvents(ctx)
}

func (u *eventUsecase) GetEventByID(ctx context.Context, id int64) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetEventByID(ctx, id)
}

func (u *eventUsecase) DeleteEvent(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.DeleteEvent(ctx, id)
}

func (u *eventUsecase) CreateLocation(ctx context.Context, req *dto.LocationRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	location := &domain.Location{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Capacity:    req.Capacity,
		// TODO update later
		OwnerID: 1,
	}

	return u.repo.CreateLocation(ctx, location)
}

func (u *eventUsecase) ListLocations(ctx context.Context, limit, offset int) ([]*domain.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	if limit == 0 {
		limit = 10
	}

	return u.repo.ListLocations(ctx, limit, offset)
}

func (u *eventUsecase) UpdateLocation(ctx context.Context, id int64, req *dto.LocationUpdateRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	location, err := u.repo.GetLocationByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		location.Name = *req.Name
	}

	if req.Description != nil {
		location.Description = *req.Description
	}

	if req.Address != nil {
		location.Address = *req.Address
	}

	if req.Capacity != nil {
		location.Capacity = *req.Capacity
	}

	if req.OwnerID != nil {
		location.OwnerID = *req.OwnerID
	}

	return u.repo.UpdateLocation(ctx, location)
}

func (u *eventUsecase) GetLocationByID(ctx context.Context, id int64) (*domain.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetLocationByID(ctx, id)
}

func (u *eventUsecase) DeleteLocation(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.DeleteLocation(ctx, id)
}
