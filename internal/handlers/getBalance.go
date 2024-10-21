package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/gofermart/internal/auth"
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/storage"
	"go.uber.org/zap"
	"time"
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func GetUserBalanceHandler(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	select {
	case <-ctx.Done():
		logger.Log.Warn("Context canceled or timeout exceeded")
		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
			"error": "Request timed out",
		})
	default:
		token := c.Cookies("jwt")

		userID := auth.GetUserID(token)

		balance, err := storage.GetUserBalance(ctx, userID)

		if err != nil {
			logger.Log.Error("Error getting user orders", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
		}

		return c.Status(fiber.StatusOK).JSON(BalanceResponse{
			Current:   balance.CurrentBalance,
			Withdrawn: balance.WithdrawnTotal,
		})
	}
}
