package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/tokenStorage"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("jwt")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	isValid := tokenStorage.CheckToken(token)

	if !isValid {
		logger.Log.Error("Token validation failed")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	return c.Next()
}
