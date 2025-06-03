package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/go-ticket-booking/internal/domain"
)

type BookingRepository interface {
	Create(ctx context.Context, tx *sql.Tx, b *domain.Booking) error
	GetByID(ctx context.Context, id int64) (*domain.Booking, error)
	ListByUserID(ctx context.Context, userID int64) ([]*domain.Booking, error)
	ListByEventID(ctx context.Context, eventID int64) ([]*domain.Booking, error)
	Update(ctx context.Context, b *domain.Booking) error
	Confirm(ctx context.Context, tx *sql.Tx, b *domain.Booking) error
	Cancel(ctx context.Context, tx *sql.Tx, b *domain.Booking) error
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

func (r *bookingRepository) GetByID(ctx context.Context, id int64) (*domain.Booking, error) {
	var booking domain.Booking

	query := `
		SELECT id, user_id, event_id, seat_id, status, created_at, confirmed_at, cancelled_at
		FROM bookings WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.EventID,
		&booking.SeatID,
		&booking.Status,
		&booking.CreatedAt,
		&booking.ConfirmedAt,
		&booking.CancelledAt,
	)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) ListByUserID(ctx context.Context, userID int64) ([]*domain.Booking, error) {
	query := `
		SELECT id, user_id, event_id, seat_id, status, created_at, confirmed_at, cancelled_at
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
			&b.CreatedAt,
			&b.ConfirmedAt,
			&b.CancelledAt,
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
		SELECT id, user_id, event_id, seat_id, status, created_at, confirmed_at, cancelled_at
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
			&b.CreatedAt,
			&b.ConfirmedAt,
			&b.CancelledAt,
		)
		if err != nil {
			return nil, err
		}
		booking = append(booking, &b)
	}

	return booking, nil
}

func (r *bookingRepository) Confirm(ctx context.Context, tx *sql.Tx, b *domain.Booking) error {
	query := `
		UPDATE bookings SET status = 'confirmed', confirmed_at = $1
		WHERE id = $2 AND status = 'pending'
	`
	res, err := tx.ExecContext(ctx, query, b.ConfirmedAt, b.ID)
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

func (r *bookingRepository) Update(ctx context.Context, b *domain.Booking) error {
	query := `
		UPDATE bookings SET 
			user_id = $1, event_id = $2, seat_id = $3 
		WHERE id = $4
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		b.UserID,
		b.EventID,
		b.SeatID,
		b.ID,
	)
	return err
}

func (r *bookingRepository) Cancel(ctx context.Context, tx *sql.Tx, b *domain.Booking) error {
	query := `
		UPDATE bookings SET status = 'cancelled', cancelled_at = $1
		WHERE id = $2 AND status IN ('pending', 'confirmed')
	`
	res, err := tx.ExecContext(ctx, query, b.CancelledAt, b.ID)
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
