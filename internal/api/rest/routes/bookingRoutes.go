package routes

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/handler"
	"github.com/codepnw/go-ticket-booking/internal/database"
	"github.com/codepnw/go-ticket-booking/internal/repository"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
)

func SetupBookingRoutes(config *rest.ConfigRestHandler) {
	app := config.App
	db := config.DB
	tx := database.NewSqlTxManager(db)

	sectRepo := repository.NewSectionRepository(db)
	seatRepo := repository.NewSeatRepository(db)
	bookRepo := repository.NewBookingRepository(db)
	uc := usecase.NewBookingUsecase(tx, bookRepo, seatRepo, sectRepo)
	handler := handler.NewBookingHandler(db, uc)

	bookRoutes := app.Group("/bookings")

	bookRoutes.Post("/", handler.CreateBooking)
	bookRoutes.Get("/:bookingID", handler.GetBookingByID)
	bookRoutes.Get("/event/:eventID", handler.GetBookingsByEvent)
	bookRoutes.Get("/user/:userID", handler.GetBookingsByUser)
	bookRoutes.Get("/", handler.GetBookingsByStatus)
	bookRoutes.Get("/seat/:seatID", handler.AvailableBooking)
	bookRoutes.Put("/:bookingID/confirm", handler.ConfirmBooking)
	bookRoutes.Put("/:bookingID/cancel", handler.CancelBooking)
	bookRoutes.Patch("/:bookingID", handler.UpdateSeat)
}
