package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/gofermart/internal/auth"
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/tokenstorage"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("jwt")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	isValid := tokenstorage.CheckToken(token)

	fmt.Println(auth.GetUserID(token))

	if !isValid {
		logger.Log.Error("Token validation failed")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	return c.Next()
}
