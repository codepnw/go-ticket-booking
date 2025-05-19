package routes

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/handler"
	"github.com/codepnw/go-ticket-booking/internal/repository"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
)

func SetupSectionRoutes(rh *rest.ConfigRestHandler) {
	app := rh.App

	repo := repository.NewSectionRepository(rh.DB)
	uc := usecase.NewSectionUsecase(repo)
	handler := handler.NewSectionHandler(uc)

	pvtRoutes := app.Group("/sections")

	pvtRoutes.Post("/", handler.CreateSection)
	pvtRoutes.Get("/", handler.ListSections)
	pvtRoutes.Get("/:id", handler.GetSection)
	pvtRoutes.Patch("/:id", handler.UpdateSection)
	pvtRoutes.Delete("/:id", handler.DeleteSection)
}
