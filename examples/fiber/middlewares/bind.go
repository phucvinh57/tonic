package middlewares

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Bind[T any](c fiber.Ctx) error {
	var req T
	if err := c.Bind().URI(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid uri", "detail": err.Error()})
	}
	if err := c.Bind().Query(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid query", "detail": err.Error()})
	}
	if err := c.Bind().Header(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid headers", "detail": err.Error()})
	}
	if c.Method() == http.MethodPost || c.Method() == http.MethodPut || c.Method() == http.MethodPatch {
		if c.HasBody() {
			if err := c.Bind().Body(&req); err != nil {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json", "detail": err.Error()})
			}
		}
	}
	if err := validate.Struct(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "detail": err.Error()})
	}
	c.Locals("data", req)
	return c.Next()
}
