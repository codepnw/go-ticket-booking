package handler

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type sectionHandler struct {
	uc        usecase.SectionUsecase
	validator *validator.Validate
}

func NewSectionHandler(uc usecase.SectionUsecase) *sectionHandler {
	return &sectionHandler{
		uc:        uc,
		validator: validator.New(),
	}
}

func (h *sectionHandler) CreateSection(ctx *fiber.Ctx) error {
	var req dto.SectionRequest

	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	section, err := h.uc.CreateSection(ctx.Context(), req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.CreatedResponse(ctx, "", section)
}

func (h *sectionHandler) ListSections(ctx *fiber.Ctx) error {
	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	sections, err := h.uc.ListSection(ctx.Context(), limit, offset)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "", sections)
}

func (h *sectionHandler) GetSection(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	section, err := h.uc.GetSection(ctx.Context(), int64(id))
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "", section)
}

func (h *sectionHandler) UpdateSection(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")
	var req dto.SectionUpdate

	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	updated, err := h.uc.UpdateSection(ctx.Context(), int64(id), req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "section updated", updated)
}

func (h *sectionHandler) DeleteSection(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	if err := h.uc.DeleteSection(ctx.Context(), int64(id)); err != nil {
		rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "section deleted", nil)
}
