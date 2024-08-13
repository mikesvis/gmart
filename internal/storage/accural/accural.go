package accural

import (
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
		CREATE TABLE IF NOT EXISTS accurals (
			id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			order_id BIGINT NOT NULL,
			sum INT NOT NULL,
			processed_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := db.Exec(createTable)

	if err != nil {
		return err
	}

	createIndex := `
		CREATE INDEX IF NOT EXISTS accurals_order_id ON orders (order_id)
	`

	_, err = db.Exec(createIndex)

	return err
}
