package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/codepnw/go-ticket-booking/internal/domain"
)

type EventRepository interface {
	Create(ctx context.Context, e *domain.Event) error
	List(ctx context.Context) ([]*domain.Event, error)
	GetByID(ctx context.Context, id int) (*domain.Event, error)
	Update(ctx context.Context, e *domain.Event) error
	Delete(ctx context.Context, id int) error
}

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, e *domain.Event) error {
	query := `
		INSERT INTO events (name, description, start_time, end_time, location_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		e.Name,
		e.Description,
		e.StartTime,
		e.EndTime,
		e.LocationID,
	).Scan(&e.ID)
}

func (r *eventRepository) List(ctx context.Context) ([]*domain.Event, error) {
	query := `
		SELECT id, name, description, start_time, end_time, location_id, created_at, updated_at
		FROM events
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.Event

	for rows.Next() {
		var e domain.Event
		err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Description,
			&e.StartTime,
			&e.EndTime,
			&e.LocationID,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &e)
	}

	return events, nil
}

func (r *eventRepository) GetByID(ctx context.Context, id int) (*domain.Event, error) {
	query := `
		SELECT id, name, description, start_time, end_time, location_id, created_at, updated_at
		FROM events WHERE id = $1
	`
	var e domain.Event

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID,
		&e.Name,
		&e.Description,
		&e.StartTime,
		&e.EndTime,
		&e.LocationID,
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	return &e, err
}

func (r *eventRepository) Update(ctx context.Context, e *domain.Event) error {
	query := `
		UPDATE events SET name = $1, description = $2, start_time = $3, end_time = $4, location_id = $5
		WHERE id = $6
	`
	res, err := r.db.ExecContext(
		ctx,
		query,
		e.Name,
		e.Description,
		e.StartTime,
		e.EndTime,
		e.LocationID,
		e.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event id %d not found", e.ID)
	}

	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event id %d not found", id)
	}

	return nil
}
