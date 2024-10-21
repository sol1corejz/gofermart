package handlers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/gofermart/internal/auth"
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/storage"
	"regexp"
	"time"
)

var luhnCheck = regexp.MustCompile(`^\d+$`)

func isValidLuhn(order string) bool {
	var sum int
	var double bool

	for i := len(order) - 1; i >= 0; i-- {
		n := int(order[i] - '0')

		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		double = !double
	}

	return sum%10 == 0
}

func CreateOrderHandler(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	select {
	case <-ctx.Done():
		logger.Log.Warn("Context canceled or timeout exceeded")
		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
			"error": "Request timed out",
		})
	default:
		orderNumber := c.Body()

		token := c.Cookies("jwt")
		userID := auth.GetUserID(token)

		if !luhnCheck.Match(orderNumber) {
			logger.Log.Error("Invalid order number")
			return c.Status(fiber.StatusUnprocessableEntity).SendString("Invalid order number")
		}

		if !isValidLuhn(string(orderNumber)) {
			logger.Log.Error("Invalid order number")
			return c.Status(fiber.StatusUnprocessableEntity).SendString("Invalid order number")
		}

		order, err := storage.GetOrderByNumber(ctx, string(orderNumber))

		if order.OrderNumber == string(orderNumber) && order.UserID == userID {
			logger.Log.Info("Order number already registered by this user")
			return c.Status(fiber.StatusOK).SendString("Order already registered by this user")
		}

		if order.OrderNumber != "" {
			logger.Log.Info("Order number already exists")
			return c.Status(fiber.StatusConflict).SendString("Order number already exists")
		}
		if err != nil {
			if errors.Is(err, storage.ErrNoSuchOrder) {
				err = storage.CreateOrder(ctx, userID.String(), string(orderNumber))
				if err != nil {
					logger.Log.Error("Error creating order")
					return c.Status(fiber.StatusInternalServerError).SendString("Error creating order")
				}
				return c.Status(fiber.StatusAccepted).SendString("Order created")
			}
			logger.Log.Error("Error checking order")
			return c.Status(fiber.StatusInternalServerError).SendString("Error checking order")
		}

		logger.Log.Error("Error creating order")
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating order")
	}
}
