package repository

import (
	"context"
	"database/sql"

	"github.com/codepnw/go-ticket-booking/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, e *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, e *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash) VALUES ($1, $2) 
		RETURNING id;
	`
	return r.db.QueryRowContext(ctx, query, e.Email, e.Password).Scan(&e.ID)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users WHERE email = $1;
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
