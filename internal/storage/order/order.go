package order

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/mikesvis/gmart/internal/domain"
	"go.uber.org/zap"
)

type Storage struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

var ErrConflict = errors.New(http.StatusText(http.StatusConflict))

func NewStorage(db *sqlx.DB, logger *zap.SugaredLogger) *Storage {
	err := bootstrap(db)
	if err != nil {
		panic(err)
	}

	return &Storage{db, logger}
}

func bootstrap(db *sqlx.DB) error {
	createTable := `
		CREATE TABLE IF NOT EXISTS orders (
			id BIGINT PRIMARY KEY,
			user_id INT NOT NULL,
			status VARCHAR(255),
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := db.Exec(createTable)

	if err != nil {
		return err
	}

	createIndex := `
		CREATE INDEX IF NOT EXISTS orders_user_id ON orders (user_id)
	`

	_, err = db.Exec(createIndex)

	return err
}

func (s *Storage) Create(ctx context.Context, orderID, userID uint64, status domain.Status) error {
	query := `INSERT INTO orders (id, user_id, status) values ($1, $2, $3) RETURNING user_id`
	_, err := s.db.ExecContext(ctx, query, orderID, userID, status)

	// не duplicate key
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code != pgerrcode.UniqueViolation {
			return err
		}
	}

	// duplicate key
	if err != nil {
		return ErrConflict
	}

	return nil
}

func (s *Storage) FindByID(ctx context.Context, orderID uint64) (*domain.Order, error) {
	var order domain.Order

	query := `SELECT id, user_id, status, created_at FROM orders WHERE id = $1`
	err := s.db.GetContext(ctx, &order, query, orderID)

	// совпадений не найдено
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	// какая-то другая ошибка
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *Storage) FindByUserID(ctx context.Context, userID uint64) ([]domain.Order, error) {
	var orders []domain.Order

	query := `
		SELECT o.id, o.user_id, o.status, COALESCE(amount, 0) AS amount, o.created_at
		FROM orders o
		LEFT JOIN accurals a ON (a.order_id = o.id AND a.amount > 0 AND o.status = $2)
		WHERE o.user_id = $1
		ORDER BY o.created_at ASC
	`

	err := s.db.SelectContext(ctx, &orders, query, userID, domain.StatusProcessed)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
