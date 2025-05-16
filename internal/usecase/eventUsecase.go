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

	// Location
	CreateLocation(ctx context.Context, req *dto.LocationRequest) error
	ListLocations(ctx context.Context, limit, offset int) ([]*domain.Location, error)
	GetLocationByID(ctx context.Context, id int) (*domain.Location, error)
	UpdateLocation(ctx context.Context, id int, req *dto.LocationUpdateRequest) error
	DeleteLocation(ctx context.Context, id int) error
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

	if err := u.repo.CreateEvent(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (u *eventUsecase) UpdateEvent(ctx context.Context, id int, req *dto.EventRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	exist, err := u.repo.GetEventByID(ctx, id)
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

	return u.repo.UpdateEvent(ctx, &event)
}

func (u *eventUsecase) ListEvents(ctx context.Context) ([]*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.ListEvents(ctx)
}

func (u *eventUsecase) GetEventByID(ctx context.Context, id int) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetEventByID(ctx, id)
}

func (u *eventUsecase) DeleteEvent(ctx context.Context, id int) error {
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

func (u *eventUsecase) UpdateLocation(ctx context.Context, id int, req *dto.LocationUpdateRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	exists, err := u.repo.GetLocationByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Name == "" {
		req.Name = exists.Name
	}

	if req.Description == "" {
		req.Description = exists.Description
	}

	if req.Address == "" {
		req.Address = exists.Address
	}

	if req.Capacity == 0 {
		req.Capacity = exists.Capacity
	}

	if req.OwnerID == 0 {
		req.OwnerID = exists.OwnerID
	}

	location := &domain.Location{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Capacity:    req.Capacity,
		// TODO update later
		OwnerID: 1,
	}

	return u.repo.UpdateLocation(ctx, location)
}

func (u *eventUsecase) GetLocationByID(ctx context.Context, id int) (*domain.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetLocationByID(ctx, id)
}

func (u *eventUsecase) DeleteLocation(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.DeleteLocation(ctx, id)
}
