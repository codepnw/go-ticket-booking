package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/codepnw/go-ticket-booking/internal/domain"
)

type EventRepository interface {
	// Event
	CreateEvent(ctx context.Context, e *domain.Event) error
	ListEvents(ctx context.Context) ([]*domain.Event, error)
	GetEventByID(ctx context.Context, id int) (*domain.Event, error)
	UpdateEvent(ctx context.Context, e *domain.Event) error
	DeleteEvent(ctx context.Context, id int) error

	// Location
	CreateLocation(ctx context.Context, l *domain.Location) error
	ListLocations(ctx context.Context, limit, offset int) ([]*domain.Location, error)
	GetLocationByID(ctx context.Context, id int) (*domain.Location, error)
	UpdateLocation(ctx context.Context, l *domain.Location) error
	DeleteLocation(ctx context.Context, id int) error
}

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(ctx context.Context, e *domain.Event) error {
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

func (r *eventRepository) ListEvents(ctx context.Context) ([]*domain.Event, error) {
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

func (r *eventRepository) GetEventByID(ctx context.Context, id int) (*domain.Event, error) {
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

func (r *eventRepository) UpdateEvent(ctx context.Context, e *domain.Event) error {
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

func (r *eventRepository) DeleteEvent(ctx context.Context, id int) error {
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

// Location
func (r *eventRepository) CreateLocation(ctx context.Context, l *domain.Location) error {
	query := `
		INSERT INTO locations (name, description, address, capacity, owner_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	return r.db.QueryRowContext(
		ctx, 
		query,
		l.Name,
		l.Description,
		l.Address,
		l.Capacity,
		l.OwnerID,
	).Scan(&l.ID)
}

func (r *eventRepository) ListLocations(ctx context.Context, limit, offset int) ([]*domain.Location, error) {
	query := `
		SELECT id, name, description, address, capacity, owner_id, created_at, updated_at
		FROM locations LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locs []*domain.Location

	for rows.Next() {
		var l domain.Location

		err := rows.Scan(
			&l.ID,
			&l.Name,
			&l.Description,
			&l.Address,
			&l.Capacity,
			&l.OwnerID,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		locs = append(locs, &l)
	}

	return locs, nil
}

func (r *eventRepository) GetLocationByID(ctx context.Context, id int) (*domain.Location, error) {
	var loc domain.Location

	query := `
		SELECT id, name, description, address, capacity, owner_id, created_at, updated_at
		FROM locations WHERE id = $1;
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&loc.ID,
		&loc.Name,
		&loc.Description,
		&loc.Address,
		&loc.Capacity,
		&loc.OwnerID,
		&loc.CreatedAt,
		&loc.UpdatedAt,
	)

	return &loc, err
}

func (r *eventRepository) UpdateLocation(ctx context.Context, l *domain.Location) error {
	query := `
		UPDATE locations SET name = $1, description = $2, address = $3, capacity = $4, owner_id = $5
		WHERE id = $6
	`
	res, err := r.db.ExecContext(
		ctx, 
		query, 
		l.Name,
		l.Description,
		l.Address,
		l.Capacity,
		l.OwnerID,
		l.ID,
	)
	if err != nil {
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return fmt.Errorf("location id %d not found", l.ID)
	}

	return nil
}

func (r *eventRepository) DeleteLocation(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM locations WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("location id %d not found", id)
	}

	return nil
}
