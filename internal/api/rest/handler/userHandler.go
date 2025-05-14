package handler

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Uc usecase.UserUsecase
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	var req dto.UserSignup
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, "")
	}

	token, err := h.Uc.CreateUser(ctx.Context(), &req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.CreatedResponse(ctx, "", token)
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var req dto.UserLogin
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, "")
	}

	token, err := h.Uc.Login(ctx.Context(), &req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "", token)
}

func (h *UserHandler) CreateProfile(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, "", nil)
}

func (h *UserHandler) GetProfile(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, "", nil)
}

func (h *UserHandler) EditProfile(ctx *fiber.Ctx) error {

	return rest.SuccessResponse(ctx, "", nil)
}
