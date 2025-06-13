package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/errs"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type BookingUsecase interface {
	Create(ctx context.Context, req *dto.CreateBookingRequest) error
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	ListByUserID(ctx context.Context, userID int64) ([]*domain.Booking, error)
	ListByEventID(ctx context.Context, eventID int64) ([]*domain.Booking, error)
	Update(ctx context.Context, bookingID int64, req *dto.UpdateBookingRequest) (*domain.Booking, error)
	UpdateStatus(ctx context.Context, tx *sql.Tx, req dto.UpdateBookingStatus) (*domain.Booking, error)
	IsAvailable(ctx context.Context, seatID int64) (bool, error)
}

type bookingUsecase struct {
	db        *sql.DB
	bookRepo  repository.BookingRepository
	eventRepo repository.EventRepository
	seatRepo  repository.SeatRepository
	sectRepo  repository.SectionRepository
}

func NewBookingUsecase(
	db *sql.DB,
	bookRepo repository.BookingRepository,
	seatRepo repository.SeatRepository,
	sectRepo repository.SectionRepository,
) BookingUsecase {
	return &bookingUsecase{
		db:       db,
		bookRepo: bookRepo,
		seatRepo: seatRepo,
		sectRepo: sectRepo,
	}
}

func (u *bookingUsecase) Create(ctx context.Context, req *dto.CreateBookingRequest) (err error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	// start tx
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// get seat
	seat, err := u.seatRepo.GetSeatByID(ctx, req.SeatID)
	if err != nil {
		return errs.ErrSeatNotFound
	}

	// get section & check event ownership
	section, err := u.sectRepo.GetByID(ctx, seat.SectionID)
	if err != nil {
		return errs.ErrSectionNotFound
	}

	if section.EventID != req.EventID {
		return errs.ErrInvalidSeatEvent
	}

	// check seat available
	isAvailable, err := u.bookRepo.IsAvailable(ctx, seat.ID)
	if err != nil {
		return err
	}

	if !isAvailable {
		return errs.ErrSeatAlreadyBooked
	}

	err = u.bookRepo.Create(ctx, tx, &domain.Booking{
		UserID:  req.UserID,
		EventID: req.EventID,
		SeatID:  req.SeatID,
		Status:  string(dto.StatusPending),
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique_booking") {
			return errors.New("user already booked")
		}
		return err
	}
	return
}

func (u *bookingUsecase) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	res, err := u.bookRepo.GetByID(ctx, id)
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

	return u.bookRepo.ListByUserID(ctx, userID)
}

func (u *bookingUsecase) ListByEventID(ctx context.Context, eventID int64) ([]*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.bookRepo.ListByEventID(ctx, eventID)
}

func (u *bookingUsecase) Update(ctx context.Context, bookingID int64, req *dto.UpdateBookingRequest) (*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	booking, err := u.bookRepo.GetByID(ctx, bookingID)
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

	if err := u.bookRepo.Update(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (u *bookingUsecase) UpdateStatus(ctx context.Context, tx *sql.Tx, req dto.UpdateBookingStatus) (*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	booking, err := u.bookRepo.GetByID(ctx, req.ID)
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
		err = u.bookRepo.Confirm(ctx, tx, booking)

	case dto.StatusCancelled:
		booking.CancelledAt = &now
		err = u.bookRepo.Cancel(ctx, tx, booking)

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

	isAvailable, err := u.bookRepo.IsAvailable(ctx, seatID)
	if err != nil {
		return false, err
	}

	if !isAvailable {
		return false, errors.New("seat is already booked")
	}

	return true, nil
}
