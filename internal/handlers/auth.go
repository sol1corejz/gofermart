package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sol1corejz/gofermart/internal/auth" // Путь к вашему auth пакету
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/storage" // Путь к вашему пакету работы с базой данных
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
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

	if existingUser.ID.String() != uuid.Nil.String() {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User already exists",
		})
	}

	token, err := auth.GenerateToken()
	if err != nil {
		logger.Log.Error("Error generating token: ", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Error hashing password: ", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	userID := uuid.New()
	auth.UserUUID = userID

	err = storage.CreateUser(userID.String(), request.Login, string(hashedPassword))
	if err != nil {
		logger.Log.Error("Error creating user: ", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(auth.TokenExp),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func LoginHandler(c *fiber.Ctx) error {
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

	if existingUser.ID.String() == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Wrong login or password",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(request.Password))
	if err != nil {
		logger.Log.Error("Error while comparing hash: ", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Wrong login or password",
		})
	}

	auth.UserUUID = existingUser.ID

	token, err := auth.GenerateToken()
	if err != nil {
		logger.Log.Error("Error generating token: ", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(auth.TokenExp),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User authorized successfully",
	})

}
