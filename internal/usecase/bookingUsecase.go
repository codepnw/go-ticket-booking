package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type BookingUsecase interface {
	Create(ctx context.Context, tx *sql.Tx, req *dto.CreateBookingRequest) error
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	ListByUserID(ctx context.Context, userID int64) ([]*domain.Booking, error)
	ListByEventID(ctx context.Context, eventID int64) ([]*domain.Booking, error)
	Update(ctx context.Context, bookingID int64, req *dto.UpdateBookingRequest) (*domain.Booking, error)
	UpdateStatus(ctx context.Context, tx *sql.Tx, req dto.UpdateBookingStatus) (*domain.Booking, error)
	IsAvailable(ctx context.Context, seatID int64) (bool, error)
}

type bookingUsecase struct {
	repo repository.BookingRepository
}

func NewBookingUsecase(repo repository.BookingRepository) BookingUsecase {
	return &bookingUsecase{repo: repo}
}

func (u *bookingUsecase) Create(ctx context.Context, tx *sql.Tx, req *dto.CreateBookingRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.Create(ctx, tx, &domain.Booking{
		UserID:  req.UserID,
		EventID: req.EventID,
		SeatID:  req.SeatID,
		Status:  string(dto.StatusPending),
	})
}

func (u *bookingUsecase) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	res, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	return res, nil
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

func (u *bookingUsecase) Update(ctx context.Context, bookingID int64, req *dto.UpdateBookingRequest) (*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	booking, err := u.repo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if req.UserID == nil && req.EventID == nil && req.SeatID == nil {
		return booking, errors.New("no fields to update")
	}

	if req.UserID != nil {
		booking.UserID = *req.UserID
	}
	if req.EventID != nil {
		booking.EventID = *req.EventID
	}
	if req.SeatID != nil {
		booking.SeatID = *req.SeatID
	}

	if err := u.repo.Update(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (u *bookingUsecase) UpdateStatus(ctx context.Context, tx *sql.Tx, req dto.UpdateBookingStatus) (*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	booking, err := u.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if booking.Status == string(req.Status) {
		return booking, errors.New("status not update")
	}

	now := time.Now()

	switch req.Status {
	case dto.StatusConfirmed:
		booking.ConfirmedAt = &now
		err = u.repo.Confirm(ctx, tx, booking)

	case dto.StatusCancelled:
		booking.CancelledAt = &now
		err = u.repo.Cancel(ctx, tx, booking)

	default:
		return nil, errors.New("invalid booking status")
	}

	if err != nil {
		return nil, err
	}

	booking.Status = string(req.Status)
	return booking, nil
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
