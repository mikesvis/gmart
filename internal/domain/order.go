package domain

import "time"

type Order struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	Status    Status    `db:"status"`
	Accural   uint64    `db:"amount"` // ВАЖНО! значения храним только в целых числах (копейки а не рубли)
	CreatedAt time.Time `db:"created_at"`
}
