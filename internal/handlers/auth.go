package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sol1corejz/gofermart/internal/auth" // Путь к вашему auth пакету
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/storage" // Путь к вашему пакету работы с базой данных
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func RegisterHandler(c *fiber.Ctx) error {

	var request RegisterRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	existingUser, err := storage.GetUserByLogin(request.Login)
	if err != nil {
		logger.Log.Error("Error while querying user: ", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	if existingUser != "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Error hashing password: ", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	userID := auth.UserUUID
	err = storage.CreateUser(userID, request.Login, string(hashedPassword))
	if err != nil {
		logger.Log.Error("Error creating user: ", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	token, err := auth.GenerateToken()
	if err != nil {
		logger.Log.Error("Error generating token: ", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
		"token":   token,
	})
}
