package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/codepnw/go-ticket-booking/internal/database"
	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/errs"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type BookingUsecase interface {
	Create(ctx context.Context, req *dto.CreateBookingRequest) error
	GetByID(ctx context.Context, id int64) (*dto.BookingResponse, error)
	ListByUserID(ctx context.Context, userID int64) ([]*dto.BookingResponse, error)
	ListByEventID(ctx context.Context, eventID int64) ([]*dto.BookingResponse, error)
	ListByStatus(ctx context.Context, status string) ([]*dto.BookingResponse, error)
	ConfirmBooking(ctx context.Context, bookingID int64) (err error)
	CancelBooking(ctx context.Context, bookingID int64) (err error)
	IsAvailable(ctx context.Context, seatID int64) (bool, error)
	UpdateSeat(ctx context.Context, bookingID, newSeatID int64) error
}

type bookingUsecase struct {
	tx       database.TxManager
	bookRepo repository.BookingRepository
	seatRepo repository.SeatRepository
	sectRepo repository.SectionRepository
}

func NewBookingUsecase(
	tx database.TxManager,
	bookRepo repository.BookingRepository,
	seatRepo repository.SeatRepository,
	sectRepo repository.SectionRepository,
) BookingUsecase {
	return &bookingUsecase{
		tx:       tx,
		bookRepo: bookRepo,
		seatRepo: seatRepo,
		sectRepo: sectRepo,
	}
}

func (u *bookingUsecase) Create(ctx context.Context, req *dto.CreateBookingRequest) (err error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.tx.WithTx(ctx, func(tx *sql.Tx) error {
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
		return nil
	})
}

func (u *bookingUsecase) GetByID(ctx context.Context, id int64) (*dto.BookingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	res, err := u.bookRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrBookingNotFound
		}
		return nil, err
	}

	return res, nil
}

func (u *bookingUsecase) ListByUserID(ctx context.Context, userID int64) ([]*dto.BookingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.bookRepo.ListByUserID(ctx, userID)
}

func (u *bookingUsecase) ListByEventID(ctx context.Context, eventID int64) ([]*dto.BookingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.bookRepo.ListByEventID(ctx, eventID)
}

func (u *bookingUsecase) ListByStatus(ctx context.Context, status string) ([]*dto.BookingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.bookRepo.ListByStatus(ctx, status)
}

func (u *bookingUsecase) ConfirmBooking(ctx context.Context, bookingID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.tx.WithTx(ctx, func(tx *sql.Tx) error {
		booking, err := u.bookRepo.GetForUpdate(ctx, tx, bookingID)
		if err != nil {
			return errs.ErrBookingNotFound
		}

		// check seat confirmed
		confirmed, err := u.bookRepo.IsSeatConfirmed(ctx, tx, booking.SeatID)
		if err != nil {
			return err
		}
		if confirmed {
			return errs.ErrSeatAlreadyBooked
		}

		// confirmed booking
		err = u.bookRepo.Confirm(ctx, tx, booking.ID)
		if err != nil {
			return err
		}

		// cancel other booking
		err = u.bookRepo.CancelOtherBooking(ctx, tx, booking.SeatID, booking.ID)
		return err
	})
}

func (u *bookingUsecase) CancelBooking(ctx context.Context, bookingID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.tx.WithTx(ctx, func(tx *sql.Tx) error {
		booking, err := u.bookRepo.GetByID(ctx, bookingID)
		if err != nil {
			return err
		}

		if booking.Status == string(dto.StatusCancelled) {
			return errs.ErrBookingAlreadyCancelled
		}
		if booking.Status == string(dto.StatusConfirmed) {
			return errs.ErrBookingAlreadyConfirmed
		}

		err = u.bookRepo.Cancel(ctx, tx, bookingID)
		return err
	})
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

func (u *bookingUsecase) UpdateSeat(ctx context.Context, bookingID, newSeatID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.tx.WithTx(ctx, func(tx *sql.Tx) error {
		booking, err := u.bookRepo.GetForUpdate(ctx, tx, bookingID)
		if err != nil {
			return errs.ErrBookingNotFound
		}

		// update status pending only
		if booking.Status != string(dto.StatusPending) {
			return errs.ErrBookingNotPending
		}

		// get seat
		seat, err := u.seatRepo.GetSeatByID(ctx, newSeatID)
		if err != nil {
			return errs.ErrSeatNotFound
		}

		// check new seat in event
		section, err := u.sectRepo.GetByID(ctx, seat.SectionID)
		if err != nil {
			return errs.ErrSectionNotFound
		}
		if section.EventID != booking.EventID {
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

		// update seat
		err = u.bookRepo.UpdateSeat(ctx, tx, booking.ID, seat.ID)
		if err != nil {
			return err
		}

		return nil
	})
}
