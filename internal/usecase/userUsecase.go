package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/errs"
	"github.com/codepnw/go-ticket-booking/internal/helper"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

const queryTimeOut = time.Second * 5

type UserUsecase interface {
	CreateUser(ctx context.Context, req *dto.UserRegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (string, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	UpdateUser(ctx context.Context, id int64, req *dto.UserUpdateRequest) (*domain.User, error)

	// Admin
	GetUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
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

func (u *userUsecase) CreateUser(ctx context.Context, req *dto.UserRegisterRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	hashed, err := u.auth.GenenrateHashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Email:     req.Email,
		Password:  hashed,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      "user",
	}

	// create user
	created, err := u.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %v", err)
	}

	// utc time -> thai time
	t, err := helper.LoadThaiTime(created.CreatedAt)
	if err != nil {
		return nil, err
	}
	created.CreatedAt = t

	return created, nil
}

func (u *userUsecase) Login(ctx context.Context, req *dto.UserLoginRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrUserNotFound
		}
		return "", err
	}

	// verify password
	if err = u.auth.VerifyPassword(req.Password, user.Password); err != nil {
		return "", err
	}

	now := time.Now().UTC()
	user.LastLoginAt = &now

	// update last login
	if err := u.userRepo.UpdateLastLogin(ctx, user); err != nil {
		return "", err
	}

	return u.auth.GenerateToken(user.ID, user.Email, user.Role)
}

func (u *userUsecase) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	if err := helper.ConvertUserTimes(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, id int64, req *dto.UserUpdateRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.Phone != nil {
		user.Phone = *req.Phone
	}

	now := time.Now().UTC()
	user.UpdatedAt = &now

	if err := u.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) GetUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	users, err := u.userRepo.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	if err := u.userRepo.DeleteUser(ctx, id); err != nil {
		return err
	}

	return nil
}
