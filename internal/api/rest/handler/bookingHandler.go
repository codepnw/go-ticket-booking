package handler

import (
	"database/sql"

	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const (
	bookingID = "bookingID"
	userID    = "userID"
	eventID   = "eventID"
	seatID    = "seatID"
)

type bookingHandler struct {
	db        *sql.DB
	uc        usecase.BookingUsecase
	validator *validator.Validate
}

func NewBookingHandler(db *sql.DB, uc usecase.BookingUsecase) *bookingHandler {
	return &bookingHandler{
		db:        db,
		uc:        uc,
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

	// begin tx
	// tx, err := h.db.BeginTx(ctx.Context(), nil)
	// if err != nil {
	// 	return rest.InternalError(ctx, err)
	// }
	// defer tx.Rollback()

	// usecase
	if err := h.uc.Create(ctx.Context(), &req); err != nil {
		return rest.InternalError(ctx, err)
	}

	// commit tx
	// if err := tx.Commit(); err != nil {
	// 	return rest.InternalError(ctx, err)
	// }

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

func (h *bookingHandler) UpdateBooking(ctx *fiber.Ctx) error {
	id, err := rest.GetParamsID(ctx, bookingID)
	if err != nil {
		return err
	}

	var req dto.UpdateBookingRequest
	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	booking, err := h.uc.Update(ctx.Context(), id, &req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "booking updated", booking)
}

func (h *bookingHandler) UpdateBookingStatus(ctx *fiber.Ctx) error {
	req := dto.UpdateBookingStatus{}

	if err := ctx.BodyParser(&req); err != nil {
		return rest.BadRequestResponse(ctx, err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		if err.Error() == "Key: 'UpdateBookingStatus.Status' Error:Field validation for 'Status' failed on the 'oneof' tag" {
			return rest.BadRequestResponse(ctx, "status: ['pending', 'confirmed', 'cancelled']")
		}
		return rest.BadRequestResponse(ctx, err.Error())
	}

	// begin tx
	tx, err := h.db.BeginTx(ctx.Context(), nil)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	defer tx.Rollback()

	// usecase
	booking, err := h.uc.UpdateStatus(ctx.Context(), tx, req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	// commit tx
	if err := tx.Commit(); err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "status updated", booking)
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
