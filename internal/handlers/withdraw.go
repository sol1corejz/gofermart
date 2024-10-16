package handlers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/gofermart/internal/auth"
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/storage"
	"go.uber.org/zap"
	"time"
)

type WithdrawRequest struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required"`
}

func WithdrawHandler(c *fiber.Ctx) error {
	var request WithdrawRequest
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
			logger.Log.Error("Error getting user balance", zap.Error(err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err = c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if balance.CurrentBalance < request.Sum {
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"error": "Insufficient funds",
			})
		}

		_, err = storage.GetOrderByNumber(ctx, request.Order)

		if err != nil {
			if errors.Is(err, storage.ErrNoSuchOrder) {
				logger.Log.Error("Error getting user order in db", zap.Error(err))
				return c.SendStatus(fiber.StatusUnprocessableEntity)
			}
			logger.Log.Error("Error getting user order", zap.Error(err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		err = storage.CreateWithdrawal(ctx, userID, request.Order, request.Sum)
		if err != nil {
			logger.Log.Error("Error creating withdrawal", zap.Error(err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

type WithdrawalsResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func GetWithdrawalsHandler(c *fiber.Ctx) error {
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

		withdrawals, err := storage.GetUserWithdrawals(ctx, userID)

		if err != nil {
			logger.Log.Error("Error getting user withdrawals", zap.Error(err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if len(withdrawals) == 0 {
			logger.Log.Info("No withdrawals found")
			return c.SendStatus(fiber.StatusNoContent)
		}

		var response []WithdrawalsResponse
		for _, withdrawal := range withdrawals {
			response = append(response, WithdrawalsResponse{
				Order:       withdrawal.OrderNumber,
				Sum:         withdrawal.Sum,
				ProcessedAt: withdrawal.ProcessedAt,
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
