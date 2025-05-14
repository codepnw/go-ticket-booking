package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(ctx *fiber.Ctx, msg string, data any) error {
	if msg == "" {
		msg = "success"
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}

func CreatedResponse(ctx *fiber.Ctx, msg string, data any) error {
	if msg == "" {
		msg = "created"
	}

	return ctx.Status(http.StatusCreated).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}

func BadRequestResponse(ctx *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "please provide valid inputs"
	}

	return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": msg})
}

func InternalError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
}
