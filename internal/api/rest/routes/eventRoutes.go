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

	pvtRoutes.Post("/", handler.AddEvent)
	pvtRoutes.Get("/", handler.ListEvents)
	pvtRoutes.Get("/:id", handler.GetEventByID)
	pvtRoutes.Patch("/:id", handler.EditEvent)
	pvtRoutes.Delete("/:id", handler.DeleteEvent)
}
