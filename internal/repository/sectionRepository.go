package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/go-ticket-booking/internal/domain"
)

type SectionRepository interface {
	Create(ctx context.Context, s *domain.Section) error
	List(ctx context.Context, limit, offset int) ([]*domain.Section, error)
	GetByID(ctx context.Context, id int64) (*domain.Section, error)
	Update(ctx context.Context, s *domain.Section) error
	Delete(ctx context.Context, id int64) error
}

type sectionRepository struct {
	db *sql.DB
}

func NewSectionRepository(db *sql.DB) SectionRepository {
	return &sectionRepository{db: db}
}

func (r *sectionRepository) Create(ctx context.Context, s *domain.Section) error {
	query := `
		INSERT INTO sections (event_id, name, seat_count)
		VALUES ($1, $2, $3) RETURNING id
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		s.EventID,
		s.Name,
		s.SeatCount,
	).Scan(&s.ID)
}

func (r *sectionRepository) List(ctx context.Context, limit, offset int) ([]*domain.Section, error) {
	query := `
		SELECT id, event_id, name, seat_count, created_at, updated_at
		FROM sections LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []*domain.Section

	for rows.Next() {
		var s domain.Section

		err := rows.Scan(
			&s.ID,
			&s.EventID,
			&s.Name,
			&s.SeatCount,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sections = append(sections, &s)
	}
	return sections, nil
}

func (r *sectionRepository) GetByID(ctx context.Context, id int64) (*domain.Section, error) {
	var sec domain.Section

	query := `
		SELECT id, event_id, name, seat_count, created_at, updated_at
		FROM sections WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sec.ID,
		&sec.EventID,
		&sec.Name,
		&sec.SeatCount,
		&sec.CreatedAt,
		&sec.UpdatedAt,
	)

	return &sec, err
}

func (r *sectionRepository) Update(ctx context.Context, s *domain.Section) error {
	query := `
		UPDATE sections SET event_id = $1, name = $2, seat_count = $3, updated_at = $4
		WHERE id = $5
	`
	res, err := r.db.ExecContext(
		ctx,
		query,
		s.EventID,
		s.Name,
		s.SeatCount,
		s.UpdatedAt,
		s.ID,
	)
	if err != nil {
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("section id not found")
	}

	return nil
}

func (r *sectionRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM section WHERE id = $1", id)
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
