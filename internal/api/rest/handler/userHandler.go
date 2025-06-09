package handler

import (
	"errors"

	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/errs"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const userCtxKey = "user"

type userHandler struct {
	uc        usecase.UserUsecase
	validator *validator.Validate
}

func NewUserHandler(uc usecase.UserUsecase) *userHandler {
	return &userHandler{
		uc: uc,
		validator: validator.New(),
	}
}

func (h *userHandler) Register(ctx *fiber.Ctx) error {
	var req dto.UserRegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	user, err := h.uc.CreateUser(ctx.Context(), &req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.CreatedResponse(ctx, "register success", dto.NewUserResponse(user))
}

func (h *userHandler) Login(ctx *fiber.Ctx) error {
	var req dto.UserLoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, errs.ErrInvalidInputData.Error())
	}

	token, err := h.uc.Login(ctx.Context(), &req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "login success", token)
}

func (h *userHandler) GetProfile(ctx *fiber.Ctx) error {
	u, ok := ctx.Locals(userCtxKey).(*domain.User)
	if !ok {
		return rest.UnauthorizedResponse(ctx)
	}

	user, err := h.uc.GetProfile(ctx.Context(), u.ID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "user profile", dto.NewUserResponse(user))
}

func (h *userHandler) UpdateProfile(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals(userCtxKey).(*domain.User)
	if !ok {
		return rest.UnauthorizedResponse(ctx)
	}

	req := dto.UserUpdateRequest{}

	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, errs.ErrInvalidInputData.Error())
	}

	res, err := h.uc.UpdateProfile(ctx.Context(), user.ID, &req)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "user updated", res)
}
