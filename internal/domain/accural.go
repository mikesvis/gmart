package domain

import (
	"time"
)

type Accrual struct {
	ID          uint64    `db:"id"`
	OrderID     uint64    `db:"order_id"`
	Status      Status    `db:"status"`
	Amount      uint64    `db:"amount"` // ВАЖНО! значения храним только в целых числах (копейки а не рубли)
	ProcessedAt time.Time `db:"processed_at"`
}
