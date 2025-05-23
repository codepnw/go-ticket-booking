package routes

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/handler"
	"github.com/codepnw/go-ticket-booking/internal/repository"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
)

func SetupBookingRoutes(rh *rest.ConfigRestHandler) {
	app := rh.App

	repo := repository.NewBookingRepository(rh.DB)
	uc := usecase.NewBookingUsecase(repo)
	handler := handler.NewBookingHandler(uc)

	bookRoutes := app.Group("/bookings")

	bookRoutes.Post("/", handler.CreateBooking)
	bookRoutes.Get("/:bookingID", handler.GetBookingByID)
	bookRoutes.Get("/event/:eventID", handler.GetBookingsByEvent)
	bookRoutes.Get("/user/:userID", handler.GetBookingsByUser)
	bookRoutes.Get("/:bookingID/confirm", handler.ConfirmBooking)
	bookRoutes.Get("/:bookingID/cancel", handler.CancelBooking)
	bookRoutes.Get("/seat/:seatID", handler.AvailableBooking)
	bookRoutes.Patch("/:bookingID", handler.UpdateBooking)
}