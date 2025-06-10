package repository

import (
	"context"
	"database/sql"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/errs"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	UpdateLastLogin(ctx context.Context, u *domain.User) error
	UpdateUser(ctx context.Context, u *domain.User) error
	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	SetUserRole(ctx context.Context, id int64, role string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (first_name, last_name, email, phone, password, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		u.Password,
		u.Role,
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return u, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, first_name, last_name, email, password, phone, role
		FROM users WHERE email = $1;
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, u *domain.User) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, u.LastLoginAt, u.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, first_name, last_name, email, phone, role, created_at, updated_at, last_login_at
		FROM users WHERE id = $1;
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users SET first_name = $1, last_name = $2, phone = $3, updated_at = $4 
		WHERE id = $5
	`
	res, err := r.db.ExecContext(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Phone,
		u.UpdatedAt,
		u.ID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, created_at, updated_at, last_login_at
		FROM users ORDER BY created_at DESC LIMIT = $1 OFFSET = $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User

		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Phone,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.LastLoginAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) SetUserRole(ctx context.Context, id int64, role string) error {
	query := `UPDATE users SET role = $1 WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, role, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}
