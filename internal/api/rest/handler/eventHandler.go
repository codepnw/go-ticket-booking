package handler

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type eventHandler struct {
	uc        usecase.EventUsecase
	validator *validator.Validate
}

func NewEventHandler(uc usecase.EventUsecase) *eventHandler {
	return &eventHandler{
		uc:        uc,
		validator: validator.New(),
	}
}

func (h *eventHandler) AddEvent(ctx *fiber.Ctx) error {
	var req dto.EventRequest
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	event, err := h.uc.CreateEvent(ctx.Context(), &req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "", event)
}

func (h *eventHandler) EditEvent(ctx *fiber.Ctx) error {
	var req dto.EventRequest
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, "")
	}

	id, _ := ctx.ParamsInt("id")

	if err := h.uc.UpdateEvent(ctx.Context(), id, &req); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "event updated", nil)
}

func (h *eventHandler) ListEvents(ctx *fiber.Ctx) error {
	events, err := h.uc.ListEvents(ctx.Context())
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "list events", events)
}

func (h *eventHandler) GetEventByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	event, err := h.uc.GetEventByID(ctx.Context(), id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "get events", event)
}

func (h *eventHandler) DeleteEvent(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	if err := h.uc.DeleteEvent(ctx.Context(), id); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "event deleted", nil)
}

func (h *eventHandler) CreateLocation(ctx *fiber.Ctx) error {
	var req dto.LocationRequest

	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.uc.CreateLocation(ctx.Context(), &req); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.CreatedResponse(ctx, "", req)
}

func (h *eventHandler) UpdateLocation(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")
	var req dto.LocationUpdateRequest

	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.uc.UpdateLocation(ctx.Context(), id, &req); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.CreatedResponse(ctx, "", req)
}

func (h *eventHandler) ListLocations(ctx *fiber.Ctx) error {
	limit := ctx.QueryInt("limit")
	offset := ctx.QueryInt("offset")

	locations, err := h.uc.ListLocations(ctx.Context(), limit, offset)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "", locations)
}

func (h *eventHandler) GetLocation(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	location, err := h.uc.GetLocationByID(ctx.Context(), id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "", location)
}

func (h *eventHandler) DeleteLocation(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	if err := h.uc.DeleteLocation(ctx.Context(), id); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "location deleted", nil)
}
