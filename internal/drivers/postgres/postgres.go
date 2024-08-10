package drivers

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mikesvis/gmart/internal/config"
)

func NewPostgres(c *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", string(c.DatabaseURI))
	if err != nil {
		panic(err)
	}

	return db, err
}
