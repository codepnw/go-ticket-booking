package routes

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/handler"
	"github.com/codepnw/go-ticket-booking/internal/repository"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
)

func SetupSeatRoutes(rh *rest.ConfigRestHandler) {
	const (
		sectionID = "/section/:sectionID"
		seatID    = "/:seatID"
		eventID   = "/event/:eventID"
	)

	app := rh.App

	repo := repository.NewSeatRepository(rh.DB)
	uc := usecase.NewSeatRepository(repo)
	handler := handler.NewSeatHandler(uc)

	seatRoutes := app.Group("/seats")

	seatRoutes.Post("/", handler.CreateSeats)
	seatRoutes.Get(sectionID, handler.GetSeatsBySectionID)
	seatRoutes.Get(eventID, handler.GetAvailableSeatsByEvent)
	seatRoutes.Patch(seatID, handler.UpdateSeat)
	seatRoutes.Delete(seatID, handler.DeleteSeat)
	seatRoutes.Delete(sectionID, handler.DeleteSeatsBySection)
}
