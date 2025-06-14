package rest

import (
	"database/sql"
	"errors"

	"github.com/codepnw/go-ticket-booking/internal/helper/auth"
	"github.com/gofiber/fiber/v2"
)

type ConfigRestHandler struct {
	App  *fiber.App
	Auth auth.Auth
	DB   *sql.DB
}

func NewRestHandler(e *ConfigRestHandler) (*ConfigRestHandler, error) {
	if e.App == nil {
		return nil, errors.New("APP is required")
	}

	if e.Auth == (auth.Auth{}) {
		return nil, errors.New("AUTH is required")
	}

	if e.DB == nil {
		return nil, errors.New("DB is required")
	}

	return e, nil
}
