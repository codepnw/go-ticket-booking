package repository

import (
	"context"
	"database/sql"
	"time"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, userId int64) error
	IsRefreshTokenValid(ctx context.Context, token string) (bool, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		SET token = $2, expires_at = $3
	`
	res, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return err
	}

	return nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return err
	}

	return nil
}

func (r *authRepository) IsRefreshTokenValid(ctx context.Context, token string) (bool, error) {
	query := `SELECT COUNT(*) FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()`

	var count int
	err := r.db.QueryRowContext(ctx, query, token).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
