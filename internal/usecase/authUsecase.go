package usecase

import (
	"context"

	"github.com/codepnw/go-ticket-booking/internal/helper/auth"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type AuthUsecase interface {
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, userID int64) error
}

type authUsecase struct {
	repo repository.AuthRepository
	auth auth.Auth
}

func NewAuthUsecase(repo repository.AuthRepository, auth auth.Auth) AuthUsecase {
	return &authUsecase{
		repo: repo,
		auth: auth,
	}
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	// verify regresh token
	user, err := u.auth.VerifyRefreshToken("Bearer " + refreshToken)
	if err != nil {
		return "", err
	}

	// check token in db
	ok, err := u.repo.IsRefreshTokenValid(ctx, refreshToken)
	if err != nil || !ok {
		return "", err
	}

	// gen new token
	accessToken, err := u.auth.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (u *authUsecase) Logout(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.DeleteRefreshToken(ctx, userID)
}