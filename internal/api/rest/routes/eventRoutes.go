package routes

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/handler"
	"github.com/codepnw/go-ticket-booking/internal/repository"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
)

func SetupEventRoutes(rh *rest.ConfigRestHandler) {
	app := rh.App

	repo := repository.NewEventRepository(rh.DB)
	uc := usecase.NewEventUsecase(repo)
	handler := handler.NewEventHandler(uc)

	pvtRoutes := app.Group("/events", rh.Auth.Authorize)
	pvtRoutes.Post("/", handler.CreateEvent)
	pvtRoutes.Get("/", handler.ListEvents)
	pvtRoutes.Get("/:id", handler.GetEventByID)
	pvtRoutes.Patch("/:id", handler.UpdateEvent)
	pvtRoutes.Delete("/:id", handler.DeleteEvent)

	// Locations
	locRoutes := app.Group("/locations")
	locRoutes.Post("/", handler.CreateLocation)
	locRoutes.Get("/", handler.ListLocations)
	locRoutes.Get("/:id", handler.GetLocation)
	locRoutes.Patch("/:id", handler.UpdateLocation)
	locRoutes.Delete("/:id", handler.DeleteLocation)
}
