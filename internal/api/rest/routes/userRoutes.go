package routes

import (
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/handler"
	"github.com/codepnw/go-ticket-booking/internal/repository"
	"github.com/codepnw/go-ticket-booking/internal/usecase"
)

func SetupUserRoutes(config *rest.ConfigRestHandler) {
	app := config.App

	authRepo := repository.NewAuthRepository(config.DB)
	authUc := usecase.NewAuthUsecase(authRepo, config.Auth)

	userRepo := repository.NewUserRepository(config.DB)
	userUc := usecase.NewUserUsecase(userRepo, authRepo, config.Auth)

	handler := handler.NewUserHandler(userUc, authUc)

	// Public Routes
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)
	app.Post("/auth/refresh-token", handler.RefreshToken)
	// TODO
	// app.Get("/auth/forgot-password", handler.ForgotPassword)
	// app.Get("/auth/reset-password", handler.ResetPassword)

	// Private Routes
	pvt := app.Group("/users", config.Auth.Authorize)
	pvt.Get("/profile", handler.GetProfile)
	pvt.Patch("/profile", handler.UpdateProfile)
	pvt.Get("/logout", handler.Logout)
	// TODO
	// pvt.Get("/change-password", handler.ChangePassword)

	// Admin
	admin := app.Group("/admin", config.Auth.Authorize)
	admin.Get("/users", handler.AdminGetUsers)
	admin.Get("/users/:id", handler.AdminGetUser)
	admin.Patch("/users/:id", handler.AdminUpdateUser)
	admin.Delete("/users/:id", handler.AdminDeleteUser)
}
