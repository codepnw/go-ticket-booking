package handler

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const (
	bookingID = "bookingID"
	userID = "userID"
	eventID = "eventID"
	seatID = "seatID"
)

type bookingHandler struct {
	uc usecase.BookingUsecase
	validator *validator.Validate
}

func NewBookingHandler(uc usecase.BookingUsecase) *bookingHandler {
	return &bookingHandler{
		uc: uc,
		validator: validator.New(),
	}
}

func (h *bookingHandler) CreateBooking(ctx *fiber.Ctx) error {
	var req dto.CreateBookingRequest

	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.uc.Create(ctx.Context(), &req); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.CreatedResponse(ctx, "booking created", req)
}

func (h *bookingHandler) GetBookingByID(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, bookingID)
	if err != nil {
		return err
	}

	booking, err := h.uc.GetByID(ctx.Context(), id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "booking by id", booking)
}

func (h *bookingHandler) GetBookingsByUser(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, userID)
	if err != nil {
		return err
	}

	bookings, err := h.uc.ListByUserID(ctx.Context(), id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "bookings by user", bookings)
}

func (h *bookingHandler) GetBookingsByEvent(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, eventID)
	if err != nil {
		return err
	}

	bookings, err := h.uc.ListByEventID(ctx.Context(), id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "bookings by event", bookings)
}

func (h *bookingHandler) ConfirmBooking(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, bookingID)
	if err != nil {
		return err
	}

	if err := h.uc.Confirm(ctx.Context(), id); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "confirmed booking", nil)
}

func (h *bookingHandler) CancelBooking(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, bookingID)
	if err != nil {
		return err
	}

	if err := h.uc.Cancel(ctx.Context(), id); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "cancelled booking", nil)
}

func (h *bookingHandler) AvailableBooking(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, seatID)
	if err != nil {
		return err
	}

	available, err := h.uc.IsAvailable(ctx.Context(), id)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "available bookings", available)
}