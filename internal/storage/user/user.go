package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Storage struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

func NewStorage(db *sqlx.DB, logger *zap.SugaredLogger) *Storage {
	err := bootstrap(db)
	if err != nil {
		panic(err)
	}

	return &Storage{db, logger}
}

func bootstrap(db *sqlx.DB) error {
	createTable := `
		CREATE TABLE IF NOT EXISTS users (
			id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			login VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`
	_, err := db.Exec(createTable)
	return err
}

func (s *Storage) ExistByLogin(ctx context.Context, login string) (bool, error) {
	var count int

	query := `SELECT COUNT(1) FROM users WHERE login = $1`
	err := s.db.QueryRowContext(ctx, query, login).Scan(&count)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

func (s *Storage) Create(ctx context.Context, login, password string) (uint64, error) {
	var userID uint64

	err := s.db.QueryRowContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", login, password).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *Storage) GetUserID(ctx context.Context, login, password string) (uint64, error) {
	var userID uint64

	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	err := s.db.QueryRowContext(ctx, query, login, password).Scan(&userID)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return userID, nil
}
