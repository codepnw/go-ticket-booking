package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
)

type SeatRepository interface {
	Create(ctx context.Context, seats []*domain.Seat) error
	GetSeatsBySectionID(ctx context.Context, sectionID int64) ([]*domain.Seat, error)
	GetAvailableSeatsByEvent(ctx context.Context, eventID int64) ([]*domain.Seat, error)
	UpdateSeat(ctx context.Context, seatID int64, input *dto.UpdateSeatRequest) error
	DeleteSeat(ctx context.Context, seatID int64) error
	DeleteSeatsBySection(ctx context.Context, sectionID int64) error
}

type seatRepository struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) Create(ctx context.Context, seats []*domain.Seat) error {
	if len(seats) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}

	query := `
		INSERT INTO seats (section_id, row_label, seat_number) 
		VALUES ($1, $2, $3)
	`
	for _, seat := range seats {
		_, err := tx.ExecContext(ctx, query, seat.SectionID, seat.RowLabel, seat.SeatNumber)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("insert seats failed: %w", err)
		}
	}

	return tx.Commit()
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

func (r *seatRepository) UpdateSeat(ctx context.Context, seatID int64, input *dto.UpdateSeatRequest) error {
	query := `UPDATE seats SET`
	params := []any{}
	i := 1

	if input.RowLabel != nil {
		query += fmt.Sprintf(" row_label = $%d,", i)
		params = append(params, *input.RowLabel)
		i++
	}

	if input.SeatNumber != nil {
		query += fmt.Sprintf(" seat_number = $%d,", i)
		params = append(params, *input.SeatNumber)
		i++
	}

	if input.IsAvailable != nil {
		query += fmt.Sprintf(" is_available = $%d,", i)
		params = append(params, *input.IsAvailable)
		i++
	}

	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id = $%d", i)
	params = append(params, seatID)

	_, err := r.db.ExecContext(ctx, query, params...)
	return err
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
		return errors.New("seat id not found")
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
		return errors.New("section id not found")
	}

	return nil
}
