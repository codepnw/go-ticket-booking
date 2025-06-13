package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/errs"
)

type SeatRepository interface {
	Create(ctx context.Context, tx *sql.Tx, seat *domain.Seat) error
	GetSeatsBySectionID(ctx context.Context, sectionID int64) ([]*domain.Seat, error)
	GetAvailableSeatsByEvent(ctx context.Context, eventID int64) ([]*domain.Seat, error)
	GetSeatByID(ctx context.Context, id int64) (*domain.Seat, error)
	UpdateSeat(ctx context.Context, s *domain.Seat) error
	DeleteSeat(ctx context.Context, seatID int64) error
	DeleteSeatsBySection(ctx context.Context, sectionID int64) error
}

type seatRepository struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) Create(ctx context.Context, tx *sql.Tx, seat *domain.Seat) error {
	query := `
		INSERT INTO seats (section_id, row_label, seat_number) 
		VALUES ($1, $2, $3)
	`
	res, err := tx.ExecContext(ctx, query, seat.SectionID, seat.RowLabel, seat.SeatNumber)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("no seat inserted")
	}

	return err
}

func (r *seatRepository) GetSeatsBySectionID(ctx context.Context, sectionID int64) ([]*domain.Seat, error) {
	query := `
		SELECT id, section_id, row_label, seat_number, is_available
		FROM seats WHERE section_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []*domain.Seat

	for rows.Next() {
		var s domain.Seat
		err := rows.Scan(
			&s.ID,
			&s.SectionID,
			&s.RowLabel,
			&s.SeatNumber,
			&s.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		seats = append(seats, &s)
	}

	return seats, nil
}

func (r *seatRepository) GetAvailableSeatsByEvent(ctx context.Context, eventID int64) ([]*domain.Seat, error) {
	query := `
		SELECT s.id, s.section_id, s.row_label, s.seat_number, s.is_available
		FROM seats s
		INNER JOIN sections sec ON s.section_id = sec.id
		WHERE sec.event_id = $1 AND s.is_available = true
	`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []*domain.Seat

	for rows.Next() {
		var s domain.Seat
		err := rows.Scan(
			&s.ID,
			&s.SectionID,
			&s.RowLabel,
			&s.SeatNumber,
			&s.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		seats = append(seats, &s)
	}

	return seats, nil
}

func (r *seatRepository) GetSeatByID(ctx context.Context, id int64) (*domain.Seat, error) {
	s := domain.Seat{}
	query := `
		SELECT id, section_id, row_label, seat_number, is_available
		FROM seats WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&s.SectionID,
		&s.RowLabel,
		&s.SectionID,
		&s.IsAvailable,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *seatRepository) UpdateSeat(ctx context.Context, s *domain.Seat) error {
	query := `
		UPDATE seats SET section_id = $1, row_label = $2, seat_number = $3, is_available = $4
		WHERE id = $5
	`
	res, err := r.db.ExecContext(
		ctx,
		query,
		s.SectionID,
		s.RowLabel,
		s.SeatNumber,
		s.IsAvailable,
		s.ID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrSeatNotFound
	}

	return nil
}

func (r *seatRepository) DeleteSeat(ctx context.Context, seatID int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM seats WHERE id = $1", seatID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrSeatNotFound
	}

	return nil
}

func (r *seatRepository) DeleteSeatsBySection(ctx context.Context, sectionID int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM seats WHERE section_id = $1", sectionID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrSeatNotFound
	}

	return nil
}
