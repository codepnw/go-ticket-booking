package routes

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/handler"
	"github.com/codepnw/go-ticket-booking/internal/repository"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
)

func SetupBookingRoutes(config *rest.ConfigRestHandler) {
	app := config.App

	repo := repository.NewBookingRepository(config.DB)
	uc := usecase.NewBookingUsecase(repo)
	handler := handler.NewBookingHandler(config.DB, uc)

	bookRoutes := app.Group("/bookings")

	bookRoutes.Post("/", handler.CreateBooking)
	bookRoutes.Post("/status", handler.UpdateBookingStatus)
	bookRoutes.Get("/:bookingID", handler.GetBookingByID)
	bookRoutes.Get("/event/:eventID", handler.GetBookingsByEvent)
	bookRoutes.Get("/user/:userID", handler.GetBookingsByUser)
	bookRoutes.Get("/seat/:seatID", handler.AvailableBooking)
	bookRoutes.Patch("/:bookingID", handler.UpdateBooking)
}
