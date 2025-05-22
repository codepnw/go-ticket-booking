package handler

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type seatHandler struct {
	uc        usecase.SeatUsecase
	validator *validator.Validate
}

func NewSeatHandler(uc usecase.SeatUsecase) *seatHandler {
	return &seatHandler{
		uc:        uc,
		validator: validator.New(),
	}
}

func (h *seatHandler) CreateSeats(ctx *fiber.Ctx) error {
	var req dto.CreateSeatsRequest

	// validate
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if len(req.Seats) == 0 {
		return rest.BadRequestResponse(ctx, "seats is empty")
	}

	// service
	if err := h.uc.CreateSeats(ctx.Context(), &req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	return rest.SuccessResponse(ctx, "insert seats success", req)
}

func (h *seatHandler) GetSeatsBySectionID(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, "sectionID")
	if err != nil {
		return err
	}

	seats, err := h.uc.GetSeatsBySectionID(ctx.Context(), id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "seats by section", seats)
}

func (h *seatHandler) GetAvailableSeatsByEvent(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, "eventID")
	if err != nil {
		return err
	}

	seats, err := h.uc.GetAvailableSeatsByEvent(ctx.Context(), id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "available seats by event", seats)
}

func (h *seatHandler) UpdateSeat(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, "seatID")
	if err != nil {
		return err
	}

	var req dto.UpdateSeatRequest

	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.uc.UpdateSeat(ctx.Context(), id, &req); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "seat updated", req)
}

func (h *seatHandler) DeleteSeat(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, "seatID")
	if err != nil {
		return err
	}

	if err := h.uc.DeleteSeat(ctx.Context(), id); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "seat deleted", nil)
}

func (h *seatHandler) DeleteSeatsBySection(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, "sectionID")
	if err != nil {
		return err
	}

	if err := h.uc.DeleteSeatsBySection(ctx.Context(), id); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "seats by section deleted", nil)
}
