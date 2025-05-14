package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/helper"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

const queryTimeOut = time.Second * 5

type UserUsecase interface {
	CreateUser(ctx context.Context, req *dto.UserSignup) (string, error)
	Login(ctx context.Context, req *dto.UserLogin) (string, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
	auth     helper.Auth
}

func NewUserUsecase(userRepo repository.UserRepository, auth helper.Auth) *userUsecase {
	return &userUsecase{
		userRepo: userRepo,
		auth:     auth,
	}
}

func (s *userUsecase) CreateUser(ctx context.Context, req *dto.UserSignup) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	hashed, err := s.auth.GenenrateHashPassword(req.Password)
	if err != nil {
		return "", err
	}

	user := domain.User{
		Email:    req.Email,
		Password: hashed,
	}

	err = s.userRepo.Create(ctx, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("create user failed: user not found")
		}
		return "", fmt.Errorf("create user failed: %v", err)
	}

	return s.auth.GenerateToken(user.ID, user.Email)
}

func (s *userUsecase) Login(ctx context.Context, req *dto.UserLogin) (string, error) {
	// find user
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	// verify password
	if err = s.auth.VerifyPassword(req.Password, user.Password); err != nil {
		return "", err
	}

	return s.auth.GenerateToken(user.ID, user.Email)
}
