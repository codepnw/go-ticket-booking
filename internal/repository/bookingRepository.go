package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/helper"
)

type BookingRepository interface {
	Create(ctx context.Context, req *domain.Booking) error
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	ListByUserID(ctx context.Context, userID int64) ([]*domain.Booking, error)
	ListByEventID(ctx context.Context, eventID int64) ([]*domain.Booking, error)
	Confirm(ctx context.Context, bookingID int64) error
	Cancel(ctx context.Context, bookingID int64) error
	IsAvailable(ctx context.Context, seatID int64) (bool, error)
}

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, req *domain.Booking) error {
	query := `
		INSERT INTO bookings (user_id, event_id, seat_id, status, booked_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query).Scan(
		&req.UserID,
		&req.EventID,
		&req.SeatID,
		&req.Status,
		&helper.LocalTime,
	)

	return err
}

func (r *bookingRepository) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	var booking domain.Booking

	query := `
		SELECT id, user_id, event_id, seat_id, status, booked_at, cancelled_at
		FROM bookings WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.EventID,
		&booking.SeatID,
		&booking.Status,
		&booking.BookedAt,
		&booking.CanceledAt,
	)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) ListByUserID(ctx context.Context, userID int64) ([]*domain.Booking, error) {
	query := `
		SELECT id, user_id, event_id, seat_id, status, booked_at, cancelled_at
		FROM bookings WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var booking []*domain.Booking

	for rows.Next() {
		var b domain.Booking
		err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.EventID,
			&b.SeatID,
			&b.Status,
			&b.BookedAt,
			&b.CanceledAt,
		)
		if err != nil {
			return nil, err
		}
		booking = append(booking, &b)
	}

	return booking, nil
}

func (r *bookingRepository) ListByEventID(ctx context.Context, eventID int64) ([]*domain.Booking, error) {
	query := `
		SELECT id, user_id, event_id, seat_id, status, booked_at, cancelled_at
		FROM bookings WHERE event_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var booking []*domain.Booking

	for rows.Next() {
		var b domain.Booking
		err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.EventID,
			&b.SeatID,
			&b.Status,
			&b.BookedAt,
			&b.CanceledAt,
		)
		if err != nil {
			return nil, err
		}
		booking = append(booking, &b)
	}

	return booking, nil
}

func (r *bookingRepository) Confirm(ctx context.Context, bookingID int64) error {
	query := `
		UPDATE bookings SET status = 'confirmed'
		WHERE id = $1 AND status = 'pending'
	`
	res, err := r.db.ExecContext(ctx, query, bookingID)
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

func (r *bookingRepository) Cancel(ctx context.Context, bookingID int64) error {
	query := `
		UPDATE bookings SET status = 'canceled', cancelled_at = $1
		WHERE id = $2
	`
	res, err := r.db.ExecContext(ctx, query, helper.LocalTime, bookingID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
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
