package accrual

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mikesvis/gmart/internal/domain"
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
		CREATE TABLE IF NOT EXISTS accruals (
			id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			order_id BIGINT NOT NULL,
			status VARCHAR(255),
			amount INT NOT NULL,
			processed_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := db.Exec(createTable)

	if err != nil {
		return err
	}

	createIndex := `
		CREATE INDEX IF NOT EXISTS accruals_order_id ON accruals (order_id)
	`

	_, err = db.Exec(createIndex)

	return err
}

func (s *Storage) GetBalanceByUserID(ctx context.Context, userID uint64) (*domain.UserBalance, error) {
	current, err := s.GetCurrentBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	withdrawn, err := s.GetWithdrawnByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := domain.UserBalance{
		Current:   current,
		Withdrawn: withdrawn,
	}

	return &result, nil
}

func (s *Storage) GetCurrentBalanceByUserID(ctx context.Context, userID uint64) (uint64, error) {
	var result uint64
	query := `
		SELECT COALESCE(SUM(a.amount), 0)
		FROM orders o
		LEFT JOIN accruals a ON (a.order_id = o.id)
		WHERE o.user_id = $1
	`

	err := s.db.QueryRowContext(ctx, query, userID).Scan(&result)

	if err != nil {
		return 0, err
	}

	return result, nil
}

func (s *Storage) GetWithdrawnByUserID(ctx context.Context, userID uint64) (uint64, error) {
	var result int64
	query := `
		SELECT COALESCE(SUM(a.amount), 0)
		FROM orders o
		LEFT JOIN accruals a ON (a.order_id = o.id AND a.amount < 0)
		WHERE o.user_id = $1
	`

	err := s.db.QueryRowContext(ctx, query, userID).Scan(&result)

	if err != nil {
		return 0, err
	}

	return uint64(result * (-1)), nil
}

func (s *Storage) CreateWithdrawn(ctx context.Context, orderID uint64, sum int64, status domain.Status) error {
	query := `INSERT INTO accruals (order_id, status, amount) values ($1, $2, $3) RETURNING id`
	_, err := s.db.ExecContext(ctx, query, orderID, status, sum)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUserWithdrawals(ctx context.Context, userID uint64) ([]domain.UserWithdrawals, error) {
	var withdrawals []domain.UserWithdrawals

	query := `
		SELECT o.id, COALESCE(ABS(a.amount), 0) as sum, a.processed_at
		FROM orders o
		INNER JOIN accruals a ON (a.order_id = o.id AND a.amount < 0 AND a.status = $2)
		WHERE o.user_id = $1
		ORDER BY a.processed_at DESC
	`

	err := s.db.SelectContext(ctx, &withdrawals, query, userID, domain.StatusProcessed)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (s *Storage) CreateAccrual(ctx context.Context, orderID uint64, status domain.Status, amount uint64) error {
	query := `INSERT INTO accruals (order_id, status, amount) values ($1, $2, $3) RETURNING id`
	_, err := s.db.ExecContext(ctx, query, orderID, status, amount)

	if err != nil {
		return err
	}

	s.logger.Infof("created accrual with %d %s %d", orderID, status, amount)

	return nil
}

func (s *Storage) UpdateOrderStatus(ctx context.Context, orderID uint64, status domain.Status) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, status, orderID)

	if err != nil {
		return err
	}

	s.logger.Infof("updated order %d with status %s", orderID, status)

	return nil
}
