package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/gofermart/cmd/config"
	"github.com/sol1corejz/gofermart/internal/handlers"
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/storage"
	"go.uber.org/zap"
)

func main() {

	config.ParseFlags()

	err := storage.Init()
	if err != nil {
		logger.Log.Error("Failed to init storage", zap.Error(err))
		return
	}

	app := fiber.New()

	app.Post("/api/user/register", handlers.RegisterHandler)

	app.Post("/api/user/login", func(c *fiber.Ctx) error {
		return nil
	})

	app.Post("/api/user/orders", func(c *fiber.Ctx) error {
		return nil
	})

	app.Get("/api/user/order", func(c *fiber.Ctx) error {
		return nil
	})

	app.Get("/api/user/balance", func(c *fiber.Ctx) error {
		return nil
	})

	app.Post("/api/user/balance/withdraw", func(c *fiber.Ctx) error {
		return nil
	})

	app.Get("/api/user/withdrawals", func(c *fiber.Ctx) error {
		return nil
	})
}
