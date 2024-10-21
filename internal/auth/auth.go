package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sol1corejz/gofermart/internal/logger"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

func GenerateToken() (string, error) {

	userID := uuid.New()

	tokenString, err := BuildJWTString(userID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func BuildJWTString(userID uuid.UUID) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},

		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string) uuid.UUID {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
	if err != nil {
		return uuid.Nil
	}

	if !token.Valid {
		logger.Log.Info("Token is not valid")
		return uuid.Nil
	}

	if claims.UserID == uuid.Nil {
		logger.Log.Warn("Parsed UserID is nil")
	}

	fmt.Println(1111111, claims)
	fmt.Println(2222222, claims.UserID)

	logger.Log.Info("Token is valid")
	return claims.UserID
}
