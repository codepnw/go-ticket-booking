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
		uc:        uc,
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

	// get user
	user, err := h.uc.GetUser(ctx.Context(), u.ID)
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

	// update user
	res, err := h.uc.UpdateUser(ctx.Context(), user.ID, &req)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "profile updated", res)
}

// ------ Admin -------
func (h *userHandler) AdminGetUsers(ctx *fiber.Ctx) error {
	// check role
	user, ok := ctx.Locals(userCtxKey).(*domain.User)
	if !ok || user.Role != string(dto.RoleAdmin) {
		return rest.ForbiddenResponse(ctx)
	}

	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	users, err := h.uc.GetUsers(ctx.Context(), limit, offset)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "list users", users)
}

func (h *userHandler) AdminGetUser(ctx *fiber.Ctx) error {
	// check role
	user, ok := ctx.Locals(userCtxKey).(*domain.User)
	if !ok || user.Role != string(dto.RoleAdmin) {
		return rest.ForbiddenResponse(ctx)
	}

	id, err := rest.GetParamsID(ctx, "id")
	if err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	u, err := h.uc.GetUser(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "get success", u)
}

func (h *userHandler) AdminUpdateUser(ctx *fiber.Ctx) error {
	// check role
	user, ok := ctx.Locals(userCtxKey).(*domain.User)
	if !ok || user.Role != string(dto.RoleAdmin) {
		return rest.ForbiddenResponse(ctx)
	}

	id, err := rest.GetParamsID(ctx, "id")
	if err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	// validate
	req := dto.UserUpdateRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, errs.ErrInvalidInputData.Error())
	}

	// update user
	updated, err := h.uc.UpdateUser(ctx.Context(), id, &req)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "user updated", updated)
}

func (h *userHandler) AdminDeleteUser(ctx *fiber.Ctx) error {
	// check role
	user, ok := ctx.Locals(userCtxKey).(*domain.User)
	if !ok || user.Role != string(dto.RoleAdmin) {
		return rest.ForbiddenResponse(ctx)
	}

	id, err := rest.GetParamsID(ctx, "id")
	if err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	// delete user
	if err := h.uc.DeleteUser(ctx.Context(), id); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "user deleted", nil)
}
