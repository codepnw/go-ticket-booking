package handler

import (
	"errors"

	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/errs"
	"github.com/codepnw/go-ticket-booking/internal/helper/auth"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userUc    usecase.UserUsecase
	authUc    usecase.AuthUsecase
	validator *validator.Validate
}

func NewUserHandler(uc usecase.UserUsecase, authUc usecase.AuthUsecase) *userHandler {
	return &userHandler{
		userUc:    uc,
		authUc:    authUc,
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

	user, err := h.userUc.CreateUser(ctx.Context(), &req)
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

	accressToken, refreshToken, err := h.userUc.Login(ctx.Context(), &req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "login success", &fiber.Map{
		"access_token":  accressToken,
		"refresh_token": refreshToken,
	})
}

func (h *userHandler) GetProfile(ctx *fiber.Ctx) error {
	u, ok := auth.GetCurrentUser(ctx)
	if !ok {
		return rest.UnauthorizedResponse(ctx)
	}

	// get user
	user, err := h.userUc.GetUser(ctx.Context(), u.ID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "user profile", dto.NewUserResponse(user))
}

func (h *userHandler) UpdateProfile(ctx *fiber.Ctx) error {
	user, ok := auth.GetCurrentUser(ctx)
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
	res, err := h.userUc.UpdateUser(ctx.Context(), user.ID, &req)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "profile updated", res)
}

// ------ Auth --------
func (h *userHandler) RefreshToken(ctx *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	refreshToken, err := h.authUc.RefreshToken(ctx.Context(), req.RefreshToken)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "token refreshed", refreshToken)
}

func (h *userHandler) Logout(ctx *fiber.Ctx) error {
	user, ok := auth.GetCurrentUser(ctx)
	if !ok {
		return rest.UnauthorizedResponse(ctx)
	}

	if err := h.authUc.Logout(ctx.Context(), user.ID); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "logout success", nil)
}

// ------ Admin -------
func (h *userHandler) AdminGetUsers(ctx *fiber.Ctx) error {
	// check role
	if ok := h.checkAdminRole(ctx); !ok {
		return rest.ForbiddenResponse(ctx)
	}

	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	users, err := h.userUc.GetUsers(ctx.Context(), limit, offset)
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
	if ok := h.checkAdminRole(ctx); !ok {
		return rest.ForbiddenResponse(ctx)
	}

	id, err := rest.GetParamsID(ctx, "id")
	if err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	u, err := h.userUc.GetUser(ctx.Context(), id)
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
	if ok := h.checkAdminRole(ctx); !ok {
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
	updated, err := h.userUc.UpdateUser(ctx.Context(), id, &req)
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
	// user, ok := auth.GetCurrentUser(ctx)
	// if !ok || user.Role != string(dto.RoleAdmin) {
	// 	return rest.ForbiddenResponse(ctx)
	// }
	if ok := h.checkAdminRole(ctx); !ok {
		return rest.ForbiddenResponse(ctx)
	}

	id, err := rest.GetParamsID(ctx, "id")
	if err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	// delete user
	if err := h.userUc.DeleteUser(ctx.Context(), id); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return rest.NotFoundResponse(ctx, err.Error())
		}
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "user deleted", nil)
}

func (h *userHandler) checkAdminRole(ctx *fiber.Ctx) bool {
	user, ok := auth.GetCurrentUser(ctx)
	if !ok || user.Role != string(dto.RoleAdmin) {
		return false
	}
	return true
}
