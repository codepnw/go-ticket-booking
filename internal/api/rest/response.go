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

func NotFoundResponse(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
		"message": msg,
	})
}

func BadRequestResponse(ctx *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "please provide valid inputs"
	}

	return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": msg})
}

func UnauthorizedResponse(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
		"message": "unauthorized",
	})
}

func ForbiddenResponse(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusForbidden).JSON(&fiber.Map{"message": "admin only"})
}

func InternalError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"error": err.Error()})
}

func GetParamsID(ctx *fiber.Ctx, key string) (int64, error) {
	id, err := ctx.ParamsInt(key)
	if err != nil {
		return 0, BadRequestResponse(ctx, "invalid params id")
	}
	return int64(id), nil
}
