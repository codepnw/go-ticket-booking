package api

import (
	"log"

	"github.com/codepnw/go-ticket-booking/config"
	"github.com/codepnw/go-ticket-booking/internal/api/rest"
	"github.com/codepnw/go-ticket-booking/internal/api/rest/routes"
	"github.com/codepnw/go-ticket-booking/internal/database"
	"github.com/codepnw/go-ticket-booking/internal/helper/auth"
	"github.com/gofiber/fiber/v2"
)

func StartServer(config config.AppConfig) {
	app := fiber.New()

	db, err := database.InitPostgresDB(config.DBAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	auth := auth.SetupAuth(config.JWTSecret, config.JWTRefreshSecret)

	rhConfig := &rest.ConfigRestHandler{
		App:  app,
		DB:   db,
		Auth: auth,
	}

	rh, err := rest.NewRestHandler(rhConfig)
	if err != nil {
		log.Fatal(err)
	}

	setupRoutes(rh)

	if err := app.Listen(config.AppPort); err != nil {
		log.Fatal(err)
	}
}

func setupRoutes(config *rest.ConfigRestHandler) {
	routes.SetupUserRoutes(config)
	routes.SetupEventRoutes(config)
	routes.SetupSectionRoutes(config)
	routes.SetupSeatRoutes(config)
	routes.SetupBookingRoutes(config)
}
