package routes

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/handler"
	"github.com/codepnw/go-ticket-booking/internal/repository"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
)

func SetupUserRoutes(rh *rest.ConfigRestHandler) {
	app := rh.App

	userRepo := repository.NewUserRepository(rh.DB)
	userUc := usecase.NewUserUsecase(userRepo, rh.Auth)

	handler := handler.UserHandler{Uc: userUc}

	// Public Routes
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	// Private Routes
	pvtRoutes := app.Group("/users", rh.Auth.Authorize)
	pvtRoutes.Post("/profile", handler.CreateProfile)
	pvtRoutes.Get("/profile", handler.GetProfile)
	pvtRoutes.Patch("/profile", handler.EditProfile)
}
