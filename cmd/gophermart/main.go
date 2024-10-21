package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sol1corejz/gofermart/cmd/config"
	"github.com/sol1corejz/gofermart/internal/handlers"
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/storage"
	"go.uber.org/zap"
)

func main() {
	config.ParseFlags()

	if err := logger.Initialize(config.LogLevel); err != nil {
		logger.Log.Fatal("Failed to initialize logger", zap.Error(err))
	}

	if err := storage.Init(); err != nil {
		logger.Log.Error("Failed to init storage", zap.Error(err))
		return
	}

	if err := run(); err != nil {
		logger.Log.Fatal("Failed to run server", zap.Error(err))
	}
}

func run() error {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,OPTIONS",
	}))

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

	logger.Log.Info("Running server", zap.String("address", config.RunAddress))
	return app.Listen(config.RunAddress)
}
