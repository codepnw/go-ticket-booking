package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/errs"
)

type BookingRepository interface {
	Create(ctx context.Context, tx *sql.Tx, b *domain.Booking) error
	GetByID(ctx context.Context, id int64) (*dto.BookingResponse, error)
	ListByUserID(ctx context.Context, userID int64) ([]*dto.BookingResponse, error)
	ListByEventID(ctx context.Context, eventID int64) ([]*dto.BookingResponse, error)
	ListByStatus(ctx context.Context, status string) ([]*dto.BookingResponse, error)
	UpdateSeat(ctx context.Context, tx *sql.Tx, bookingID, seatID int64) error
	GetForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*domain.Booking, error)
	IsSeatConfirmed(ctx context.Context, tx *sql.Tx, seatID int64) (bool, error)
	Confirm(ctx context.Context, tx *sql.Tx, bookingID int64) error
	Cancel(ctx context.Context, tx *sql.Tx, bookingID int64) error
	CancelOtherBooking(ctx context.Context, tx *sql.Tx, seatID, bookingID int64) error
	IsAvailable(ctx context.Context, seatID int64) (bool, error)
}

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, tx *sql.Tx, b *domain.Booking) error {
	query := `
		INSERT INTO bookings (user_id, event_id, seat_id, status)
		VALUES ($1, $2, $3, $4) RETURNING id
	`
	err := tx.QueryRowContext(
		ctx,
		query,
		&b.UserID,
		&b.EventID,
		&b.SeatID,
		&b.Status,
	).Scan(&b.ID)

	return err
}

var selectQuery = `
	SELECT b.id, b.user_id, u.first_name, u.last_name, u.email, b.event_id, e.name, b.seat_id,
			s.row_label, s.seat_number, b.status,
			b.created_at, b.confirmed_at, b.cancelled_at
	FROM bookings b
	JOIN events e ON b.event_id = e.id
	JOIN seats s ON b.seat_id = s.id
	JOIN users u ON b.user_id = u.id
	WHERE b.
`

func (r *bookingRepository) GetByID(ctx context.Context, id int64) (*dto.BookingResponse, error) {
	var res dto.BookingResponse

	query := selectQuery + "id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID,
		&res.User.UserID,
		&res.User.FirstName,
		&res.User.LastName,
		&res.User.Email,
		&res.Event.EventID,
		&res.Event.EventName,
		&res.Seat.SeatID,
		&res.Seat.RowLabel,
		&res.Seat.SeatNumber,
		&res.Status,
		&res.CreatedAt,
		&res.ConfirmedAt,
		&res.CancelledAt,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *bookingRepository) ListByUserID(ctx context.Context, userID int64) ([]*dto.BookingResponse, error) {
	bookings, err := r.listBookings(ctx, "user_id", userID)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) ListByEventID(ctx context.Context, eventID int64) ([]*dto.BookingResponse, error) {
	bookings, err := r.listBookings(ctx, "event_id", eventID)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) ListByStatus(ctx context.Context, status string) ([]*dto.BookingResponse, error) {
	bookings, err := r.listBookings(ctx, "status", status)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) GetForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*domain.Booking, error) {
	query := `
		SELECT id, user_id, event_id, seat_id, status FROM bookings
		WHERE id = $1 FOR UPDATE
	`
	var b domain.Booking
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&b.ID,
		&b.UserID,
		&b.EventID,
		&b.SeatID,
		&b.Status,
	)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *bookingRepository) IsSeatConfirmed(ctx context.Context, tx *sql.Tx, seatID int64) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM bookings 
		WHERE seat_id = $1 AND status = 'confirmed'
	`
	var count int
	err := tx.QueryRowContext(ctx, query, seatID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *bookingRepository) Confirm(ctx context.Context, tx *sql.Tx, bookingID int64) error {
	query := `
		UPDATE bookings SET status = 'confirmed', confirmed_at = NOW()
		WHERE id = $1 AND status = 'pending'
	`
	res, err := tx.ExecContext(ctx, query, bookingID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("booking not found or already confirmed")
	}

	return nil
}

func (r *bookingRepository) Cancel(ctx context.Context, tx *sql.Tx, bookingID int64) error {
	query := `
		UPDATE bookings SET status = 'cancelled', cancelled_at = NOW()
		WHERE id = $1
	`
	res, err := tx.ExecContext(ctx, query, bookingID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("booking not found or already confirmed")
	}

	return nil
}

func (r *bookingRepository) UpdateSeat(ctx context.Context, tx *sql.Tx, bookingID, seatID int64) error {
	query := `
		UPDATE bookings SET seat_id = $1, updated_at = NOW() 
		WHERE id = $2
	`
	res, err := tx.ExecContext(ctx, query, seatID, bookingID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrBookingNotFound
	}

	return nil
}

func (r *bookingRepository) CancelOtherBooking(ctx context.Context, tx *sql.Tx, seatID, bookingID int64) error {
	query := `
		UPDATE bookings 
		SET status = 'cancelled', cancelled_at = NOW() 
		WHERE seat_id = $1 AND status = 'pending' AND id != $2
	`
	_, err := tx.ExecContext(ctx, query, seatID, bookingID)
	return err
}

func (r *bookingRepository) IsAvailable(ctx context.Context, seatID int64) (bool, error) {
	var count int

	query := `
		SELECT COUNT(*) FROM bookings
		WHERE seat_id = $1 AND status = 'confirmed'
	`
	err := r.db.QueryRowContext(ctx, query, seatID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check available: %w", err)
	}
	return count == 0, nil
}

func (r *bookingRepository) listBookings(ctx context.Context, where string, data any) ([]*dto.BookingResponse, error) {
	query := fmt.Sprintf("%s%s%s", selectQuery, where, "=$1")

	rows, err := r.db.QueryContext(ctx, query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*dto.BookingResponse

	for rows.Next() {
		var res dto.BookingResponse
		err := rows.Scan(
			&res.ID,
			&res.User.UserID,
			&res.User.FirstName,
			&res.User.LastName,
			&res.User.Email,
			&res.Event.EventID,
			&res.Event.EventName,
			&res.Seat.SeatID,
			&res.Seat.RowLabel,
			&res.Seat.SeatNumber,
			&res.Status,
			&res.CreatedAt,
			&res.ConfirmedAt,
			&res.CancelledAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &res)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
