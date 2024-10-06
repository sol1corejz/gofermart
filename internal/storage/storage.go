package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sol1corejz/gofermart/cmd/config"
	"github.com/sol1corejz/gofermart/internal/logger"
	"github.com/sol1corejz/gofermart/internal/models"
	"go.uber.org/zap"
)

var (
	DB                     *sql.DB
	ErrConnectionFailed    = errors.New("db connection failed")
	ErrCreatingTableFailed = errors.New("creating table failed")
)

func Init() error {
	if config.DatabaseURI == "" {
		return ErrConnectionFailed
	}

	db, err := sql.Open("pgx", config.DatabaseURI)
	if err != nil {
		logger.Log.Fatal("Error opening database connection", zap.Error(err))
		return ErrConnectionFailed
	}
	DB = db

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			login VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id),
			order_number VARCHAR(255) UNIQUE NOT NULL,
			status VARCHAR(20) NOT NULL,
			accrual DECIMAL(10, 2),
			uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS user_balances (
			user_id UUID NOT NULL REFERENCES users(id),
			current_balance DECIMAL(10, 2) NOT NULL DEFAULT 0,
			withdrawn_total DECIMAL(10, 2) NOT NULL DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS withdrawals (
			id SERIAL PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id),
			order_number VARCHAR(255) NOT NULL,
			sum DECIMAL(10, 2) NOT NULL,
			processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			logger.Log.Error("Error creating table", zap.Error(err))
			return ErrCreatingTableFailed
		}
	}

	return nil
}

func GetUserByLogin(login string) (models.User, error) {

	var existingUser models.User

	err := DB.QueryRow(`
		SELECT * FROM users WHERE login = $1;
	`, login).Scan(&existingUser.ID, &existingUser.Login, &existingUser.PasswordHash, &existingUser.CreatedAt)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.User{}, err
		}
	}

	return existingUser, nil
}

func CreateUser(userID string, login string, passwordHash string) error {
	fmt.Printf("Creating user %s with login %s with password %s\n", userID, login, passwordHash)
	_, err := DB.Exec(`
		INSERT INTO users (id, login, password_hash) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING;
	`, userID, login, passwordHash)

	if err != nil {
		return err
	}

	return nil
}
