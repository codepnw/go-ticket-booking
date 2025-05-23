package usecase

import (
	"context"
	"errors"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type BookingUsecase interface {
	Create(ctx context.Context, req *dto.CreateBookingRequest) error
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	ListByUserID(ctx context.Context, userID int64) ([]*domain.Booking, error)
	ListByEventID(ctx context.Context, eventID int64) ([]*domain.Booking, error)
	Update(ctx context.Context, bookingID int64, input *dto.UpdateBookingRequest) error
	Confirm(ctx context.Context, bookingID int64) error
	Cancel(ctx context.Context, bookingID int64) error
	IsAvailable(ctx context.Context, seatID int64) (bool, error)
}

type bookingUsecase struct {
	repo repository.BookingRepository
}

func NewBookingUsecase(repo repository.BookingRepository) BookingUsecase {
	return &bookingUsecase{repo: repo}
}

func (u *bookingUsecase) Create(ctx context.Context, req *dto.CreateBookingRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.Create(ctx, &domain.Booking{
		UserID:  req.UserID,
		EventID: req.EventID,
		SeatID:  req.SeatID,
		Status:  string(dto.StatusPending),
	})
}

func (u *bookingUsecase) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *bookingUsecase) ListByUserID(ctx context.Context, userID int64) ([]*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.ListByUserID(ctx, userID)
}

func (u *bookingUsecase) ListByEventID(ctx context.Context, eventID int64) ([]*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.ListByEventID(ctx, eventID)
}

func (u *bookingUsecase) Update(ctx context.Context, bookingID int64, input *dto.UpdateBookingRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	if input.UserID == nil && input.EventID == nil && input.SeatID == nil {
		return errors.New("no fields to update")
	}

	return u.repo.Update(ctx, bookingID, input)
}

func (u *bookingUsecase) Confirm(ctx context.Context, bookingID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.Confirm(ctx, bookingID)
}

func (u *bookingUsecase) Cancel(ctx context.Context, bookingID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.Cancel(ctx, bookingID)
}

func (u *bookingUsecase) IsAvailable(ctx context.Context, seatID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	isAvailable, err := u.repo.IsAvailable(ctx, seatID)
	if err != nil {
		return false, err
	}

	if !isAvailable {
		return false, errors.New("seat is already booked")
	}

	return true, nil
}
