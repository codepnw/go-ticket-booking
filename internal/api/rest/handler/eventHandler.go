package handler

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type eventHandler struct {
	uc usecase.EventUsecase
}

func NewEventHandler(uc usecase.EventUsecase) *eventHandler {
	return &eventHandler{uc: uc}
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