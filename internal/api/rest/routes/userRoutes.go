package routes

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/handler"
	"github.com/codepnw/go-ticket-booking/internal/repository"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
)

func SetupUserRoutes(config *rest.ConfigRestHandler) {
	app := config.App

	userRepo := repository.NewUserRepository(config.DB)
	userUc := usecase.NewUserUsecase(userRepo, config.Auth)
	handler := handler.NewUserHandler(userUc)

	// Public Routes
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)
	// TODO
	// app.Get("/auth/refresh-token", handler.RefreshToken)
	// app.Get("/auth/forgot-password", handler.ForgotPassword)
	// app.Get("/auth/reset-password", handler.ResetPassword)

	// Private Routes
	pvt := app.Group("/users", config.Auth.Authorize)
	pvt.Get("/profile", handler.GetProfile)
	pvt.Patch("/profile", handler.UpdateProfile)
	// TODO
	// pvt.Get("/logout", handler.Logout)
	// pvt.Get("/change-password", handler.ChangePassword)

	// Admin
	admin := app.Group("/admin", config.Auth.Authorize)
	admin.Get("/users", handler.AdminGetUsers)
	admin.Get("/users/:id", handler.AdminGetUser)
	admin.Patch("/users/:id", handler.AdminUpdateUser)
	admin.Delete("/users/:id", handler.AdminDeleteUser)
}
